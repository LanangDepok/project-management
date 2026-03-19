package controllers

import (
	"time"

	"github.com/LanangDepok/project-management/models"
	"github.com/LanangDepok/project-management/services"
	"github.com/LanangDepok/project-management/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CardController struct {
	services services.CardService
}

func NewCardController(s services.CardService) CardController {
	return CardController{services: s}
}

func (c *CardController) CreateCard(ctx fiber.Ctx) error {
	type CreateCardRequest struct {
		ListPublicID string    `json:"list_id"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		DueDate      time.Time `json:"due_date"`
		Position     int       `json:"position"`
	}

	var req CreateCardRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing request body", err.Error())
	}

	card := &models.Card{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     &req.DueDate,
		Position:    req.Position,
	}

	if err := c.services.Create(card, req.ListPublicID); err != nil {
		return utils.InternalServerError(ctx, "Gagal membuat card", err.Error())
	}

	return utils.Success(ctx, "Card berhasil dibuat", nil)
}

func (c *CardController) UpdateCard(ctx fiber.Ctx) error {
	publicID := ctx.Params("id")

	type updateCardRequest struct {
		ListPublicID string     `json:"list_id"`
		Title        string     `json:"title"`
		Description  string     `json:"description"`
		DueDate      *time.Time `json:"due_date"`
		Position     int        `json:"position"`
	}

	var req updateCardRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing request body", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	card := &models.Card{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Position:    req.Position,
		PublicID:    uuid.MustParse(publicID),
	}

	if err := c.services.Update(card, req.ListPublicID); err != nil {
		return utils.InternalServerError(ctx, "Gagal memperbarui card", err.Error())
	}

	return utils.Success(ctx, "Card berhasil diperbarui", nil)
}

func (c *CardController) DeleteCard(ctx fiber.Ctx) error {
	publicID := ctx.Params("id")

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	card, err := c.services.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Card tidak ditemukan", err.Error())
	}

	if err := c.services.Delete(uint(card.InternalID)); err != nil {
		return utils.BadRequest(ctx, "Gagal menghapus data", err.Error())
	}

	return utils.Success(ctx, "Card berhasil dihapus", nil)
}

func (c *CardController) GetCardDetail(ctx fiber.Ctx) error {
	cardPublicID := ctx.Params("id")

	card, err := c.services.GetByPublicID(cardPublicID)
	if err != nil {
		return utils.NotFound(ctx, "Card tidak ditemukan", err.Error())
	}

	return utils.Success(ctx, "Card berhasil ditemukan", card)
}
