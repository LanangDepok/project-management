package services

import (
	"errors"
	"fmt"

	"github.com/LanangDepok/project-management/config"
	"github.com/LanangDepok/project-management/models"
	"github.com/LanangDepok/project-management/models/types"
	"github.com/LanangDepok/project-management/repositories"
	"github.com/LanangDepok/project-management/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type listService struct {
	listRepo    repositories.ListRepository
	boardRepo   repositories.BoardRepository
	listPosRepo repositories.ListPositionRepository
}

type ListWithOrder struct {
	Position []uuid.UUID
	Lists    []models.List
}

type ListService interface {
	GetByBoardID(boardPublicID string) (*ListWithOrder, error)
	GetByID(id uint) (*models.List, error)
	GetByPublicID(publicID string) (*models.List, error)
	Create(list *models.List) error
	Update(list *models.List) error
	Delete(id uint) error
	UpdatePosition(boardPublicID string, positions []uuid.UUID) error
}

func NewListService(listRepo repositories.ListRepository, boardRepo repositories.BoardRepository, listPosRepo repositories.ListPositionRepository) ListService {
	return &listService{listRepo, boardRepo, listPosRepo}
}

func (s *listService) GetByBoardID(boardPublicID string) (*ListWithOrder, error) {
	//verifikasi board
	_, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return nil, errors.New("board not found")
	}

	position, err := s.listPosRepo.GetListOrder(boardPublicID)
	if err != nil {
		return nil, errors.New("failed to get list order: " + err.Error())
	}

	lists, err := s.listRepo.FindByBoardID(boardPublicID)
	if err != nil {
		return nil, errors.New("failed to get list: " + err.Error())
	}

	//sorting by position
	orderedList := utils.SortListsByPosition(lists, position)

	return &ListWithOrder{
		Position: position,
		Lists:    orderedList,
	}, nil
}

func (s *listService) GetByID(id uint) (*models.List, error) {
	return s.listRepo.FindByID(id)
}

func (s *listService) GetByPublicID(publicID string) (*models.List, error) {
	return s.listRepo.FindByPublicID(publicID)
}

func (s *listService) Create(list *models.List) error {
	//validasi board
	board, err := s.boardRepo.FindByPublicID(list.BoardPublicID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("board not found")
		}
		return fmt.Errorf("failed to get board: %v", err)
	}
	list.BoardInternalID = board.InternalID

	if list.PublicID == uuid.Nil {
		list.PublicID = uuid.New()
	}

	//transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//simpan list baru
	if err := tx.Create(list).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create list: %v", err)
	}

	//update position
	var position models.ListPosition
	res := tx.Where("board_internal_id = ?", board.InternalID).First(&position)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		//buat baru jika belum ada
		position = models.ListPosition{
			PublicID:  uuid.New(),
			BoardID:   board.InternalID,
			ListOrder: types.UUIDArray{list.PublicID},
		}

		if err := tx.Create(&position).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create list position: %v", err)
		}
	} else if res.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create list position: %v", res.Error)
	} else {
		//tambahkan ID baru
		position.ListOrder = append(position.ListOrder, list.PublicID)

		//update ke DB
		if err := tx.Model(&position).Update("list_order", position.ListOrder).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update list position: %v", err)
		}
	}

	//commit trx
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (s *listService) Update(list *models.List) error {
	return s.listRepo.Update(list)
}

func (s *listService) Delete(id uint) error {
	return s.listRepo.Delete(id)
}

func (s *listService) UpdatePosition(boardPublicID string, positions []uuid.UUID) error {
	//verisifikasi board
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil {
		return errors.New("board not found")
	}

	//get list position
	position, err := s.listPosRepo.GetByBoard(board.PublicID.String())
	if err != nil {
		return errors.New("failed to get list position: " + err.Error())
	}

	//update list order
	position.ListOrder = positions
	return s.listPosRepo.UpdateListOrder(position)
}
