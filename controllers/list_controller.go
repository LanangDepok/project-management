package controllers

import (
	"github.com/LanangDepok/project-management/models"
	"github.com/LanangDepok/project-management/services"
	"github.com/LanangDepok/project-management/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ListController struct {
	service services.ListService
}

func NewListController(s services.ListService) *ListController {
	return &ListController{service: s}
}

// CreateList godoc
// @Summary      Create a list
// @Description  Create a new list
// @Tags         lists
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      models.List  true  "List data"
// @Success      200   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      500   {object}  utils.Response
// @Router       /api/v1/lists [post]
func (c *ListController) CreateList(ctx fiber.Ctx) error {
	list := new(models.List)

	if err := ctx.Bind().JSON(list); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing request body", err.Error())
	}

	if err := c.service.Create(list); err != nil {
		return utils.BadRequest(ctx, "Gagal membuat list", err.Error())
	}

	return utils.Success(ctx, "List berhasil dibuat", nil)
}

// UpdateList godoc
// @Summary      Update a list
// @Description  Update an existing list by public UUID
// @Tags         lists
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string  true  "List Public UUID"
// @Param        body  body      models.List  true  "List data"
// @Success      200   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      404   {object}  utils.Response
// @Failure      500   {object}  utils.Response
// @Router       /api/v1/lists/{id} [put]
func (c *ListController) UpdateList(ctx fiber.Ctx) error {
	publicID := ctx.Params("id")
	list := new(models.List)

	if err := ctx.Bind().JSON(list); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing request body", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	existingList, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.BadRequest(ctx, "List tidak ditemukan", err.Error())
	}

	list.InternalID = existingList.InternalID
	list.PublicID = existingList.PublicID

	if err := c.service.Update(list); err != nil {
		return utils.BadRequest(ctx, "Gagal memperbarui list", err.Error())
	}

	updatedList, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan", err.Error())
	}

	return utils.Success(ctx, "List berhasil diperbarui", updatedList)
}

// GetListOnBoard godoc
// @Summary      Get all lists on board
// @Description  Get all lists on board
// @Tags         lists
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        board_id    path      string  true  "Board Public UUID"
// @Success      200   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      404   {object}  utils.Response
// @Failure      500   {object}  utils.Response
// @Router       /api/v1/boards/{board_id}/lists [get]
func (c *ListController) GetListOnBoard(ctx fiber.Ctx) error {
	boardPublicID := ctx.Params("board_id")
	if _, err := uuid.Parse(boardPublicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	lists, err := c.service.GetByBoardID(boardPublicID)
	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan", err.Error())
	}

	return utils.Success(ctx, "List berhasil ditemukan", lists)
}

func (c *ListController) DeleteList(ctx fiber.Ctx) error {
	publicID := ctx.Params("id")
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	list, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List tidak ditemukan", err.Error())
	}

	if err := c.service.Delete(uint(list.InternalID)); err != nil {
		return utils.InternalServerError(ctx, "Gagal menghapus list", err.Error())
	}

	return utils.Success(ctx, "List berhasil dihapus", nil)
}
