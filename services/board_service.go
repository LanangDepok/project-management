package services

import (
	"errors"

	"github.com/LanangDepok/project-management/models"
	"github.com/LanangDepok/project-management/repositories"
	"github.com/google/uuid"
)

type BoardService interface {
	Create(board *models.Board) error
	Update(board *models.Board) error
	GetByPublicID(publicID string) (*models.Board, error)
}

type boardService struct {
	boardRepo repositories.BoardRepository
	userRepo  repositories.UserRepository
}

func NewBoardService(boardRepo repositories.BoardRepository, userRepo repositories.UserRepository) BoardService {
	return &boardService{boardRepo, userRepo}
}

func (s *boardService) Create(board *models.Board) error {
	owner, err := s.userRepo.FindByPublicID(board.OwnerPublicID.String())
	if err != nil {
		return errors.New("owner tidak ditemukan")
	}
	board.PublicID = uuid.New()
	board.OwnerID = owner.InternalID
	return s.boardRepo.Create(board)
}

func (s *boardService) Update(board *models.Board) error {
	return s.boardRepo.Update(board)
}

func (s *boardService) GetByPublicID(publicID string) (*models.Board, error) {
	return s.boardRepo.FindByPublicID(publicID)
}
