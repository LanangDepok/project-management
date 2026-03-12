package models

import (
	"time"

	"github.com/google/uuid"
)

type Board struct {
	InternalID    int64      `json:"internal_id"     gorm:"primaryKey;autoIncrement"`
	PublicID      uuid.UUID  `json:"public_id"       gorm:"column:public_id;type:uuid"`
	Title         string     `json:"title"           gorm:"column:title"`
	Description   string     `json:"description"     gorm:"column:description"`
	OwnerID       int64      `json:"owner_internal_id" gorm:"column:owner_internal_id"`
	OwnerPublicID uuid.UUID  `json:"owner_public_id" gorm:"column:owner_public_id;type:uuid"`
	CreatedAt     time.Time  `json:"created_at"`
	DueDate       *time.Time `json:"due_date,omitempty" gorm:"column:due_date"`
}

// CreateBoardRequest is the body for POST /api/v1/boards.
type CreateBoardRequest struct {
	Title       string     `json:"title"       example:"Sprint Board"`
	Description string     `json:"description" example:"Q3 2025 sprint tasks"`
	DueDate     *time.Time `json:"due_date,omitempty" example:"2025-12-31T00:00:00Z"`
}

type UpdateBoardRequest struct {
	Title       string     `json:"title"       example:"Sprint Board Updated"`
	Description string     `json:"description" example:"Updated description"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

type AddBoardMembersRequest struct {
	UserIDs []string `json:"user_ids" example:"uuid1,uuid2"`
}

type RemoveBoardMembersRequest struct {
	UserIDs []string `json:"user_ids" example:"uuid1,uuid2"`
}
