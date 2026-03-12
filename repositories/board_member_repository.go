package repositories

import (
	"github.com/LanangDepok/project-management/config"
	"github.com/LanangDepok/project-management/models"
)

type BoardMemberRepository interface {
	GetMembers(boardPublicID string) ([]models.User, error)
}

type boardMemberRepository struct {
}

func NewBoardMemberRepository() BoardMemberRepository {
	return &boardMemberRepository{}
}

func (r *boardMemberRepository) GetMembers(boardPublicID string) ([]models.User, error) {
	var users []models.User
	err := config.DB.Joins("LEFT JOIN board_members ON board_members.user_internal_id = users.internal_id").
		Joins("JOIN boards on board_members.board_internal_id = boards.internal_id").
		Where("boards.public_id = ?", boardPublicID).
		Find(&users).Error
	return users, err
}
