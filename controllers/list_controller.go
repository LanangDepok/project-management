package controllers

import (
	"github.com/LanangDepok/project-management/models"
	"github.com/LanangDepok/project-management/services"
	"github.com/LanangDepok/project-management/utils"
	"github.com/gofiber/fiber/v3"
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
