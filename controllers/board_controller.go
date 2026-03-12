package controllers

import (
	"math"
	"strconv"

	"github.com/LanangDepok/project-management/models"
	"github.com/LanangDepok/project-management/services"
	"github.com/LanangDepok/project-management/utils"
	"github.com/gofiber/fiber/v3"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type BoardController struct {
	service services.BoardService
}

func NewBoardController(s services.BoardService) *BoardController {
	return &BoardController{service: s}
}

// CreateBoard godoc
// @Summary      Create a board
// @Description  Create a new project board for the authenticated user
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      models.CreateBoardRequest  true  "Board payload"
// @Success      201   {object}  utils.Response{data=models.Board}
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Router       /api/v1/boards [post]
func (c *BoardController) CreateBoard(ctx fiber.Ctx) error {
	// Extract JWT claims
	token, ok := ctx.Locals("user").(*jwt.Token)
	if !ok {
		return utils.Unauthorized(ctx, "Token tidak valid", "invalid token type")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return utils.Unauthorized(ctx, "Token tidak valid", "invalid claims type")
	}

	pubIDStr, ok := claims["pub_id"].(string)
	if !ok {
		return utils.BadRequest(ctx, "Klaim token tidak valid", "pub_id missing")
	}
	ownerPublicID, err := uuid.Parse(pubIDStr)
	if err != nil {
		return utils.BadRequest(ctx, "Gagal parse user ID", err.Error())
	}

	var req models.CreateBoardRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing request body", err.Error())
	}

	board := &models.Board{
		Title:         req.Title,
		Description:   req.Description,
		DueDate:       req.DueDate,
		OwnerPublicID: ownerPublicID,
	}

	if err := c.service.Create(board); err != nil {
		return utils.BadRequest(ctx, "Gagal membuat board", err.Error())
	}
	return utils.Created(ctx, "Board berhasil dibuat", board)
}

// UpdateBoard godoc
// @Summary      Update a board
// @Description  Update board data by public UUID
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                    true  "Board Public UUID"
// @Param        body  body      models.UpdateBoardRequest true  "Update payload"
// @Success      200   {object}  utils.Response{data=models.Board}
// @Failure      400   {object}  utils.Response
// @Failure      404   {object}  utils.Response
// @Router       /api/v1/boards/{id}/members [put]
func (c *BoardController) UpdateBoard(ctx fiber.Ctx) error {
	publicID := ctx.Params("id")

	var req models.UpdateBoardRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing request body", err.Error())
	}
	board := &models.Board{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
	}

	if err := ctx.Bind().JSON(board); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing request body", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	existingBoard, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "Board tidak ditemukan", err.Error())
	}

	board.InternalID = existingBoard.InternalID
	board.PublicID = existingBoard.PublicID
	board.OwnerID = existingBoard.OwnerID
	board.OwnerPublicID = existingBoard.OwnerPublicID
	board.CreatedAt = existingBoard.CreatedAt

	if err := c.service.Update(board); err != nil {
		return utils.BadRequest(ctx, "Gagal memperbarui board", err.Error())
	}
	return utils.Success(ctx, "Board berhasil diperbarui", board)
}

// AddBoardMembers godoc
// @Summary      Add members to a board
// @Description  Add members to a board by public UUID
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                         true  "Board Public UUID"
// @Param        body  body      models.AddBoardMembersRequest  true  "List of user public IDs"
// @Success      200   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      404   {object}  utils.Response
// @Router       /api/v1/boards/{id}/members [post]
func (c *BoardController) AddBoardMembers(ctx fiber.Ctx) error {
	publicID := ctx.Params("id")

	var req models.AddBoardMembersRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing request body", err.Error())
	}

	if len(req.UserIDs) == 0 {
		return utils.BadRequest(ctx, "User IDs tidak boleh kosong", "user_ids is required")
	}

	if err := c.service.AddMember(publicID, req.UserIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal menambahkan member", err.Error())
	}

	return utils.Success(ctx, "Member berhasil ditambahkan", nil)
}

// RemoveBoardMembers godoc
// @Summary      Remove members from a board
// @Description  Remove members from a board by public UUID
// @Tags         boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                         true  "Board Public UUID"
// @Param        body  body      models.RemoveBoardMembersRequest  true  "List of user public IDs"
// @Success      200   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      404   {object}  utils.Response
// @Router       /api/v1/boards/{id}/members [delete]
func (c *BoardController) RemoveBoardMembers(ctx fiber.Ctx) error {
	publicID := ctx.Params("id")

	var req models.RemoveBoardMembersRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing request body", err.Error())
	}

	if len(req.UserIDs) == 0 {
		return utils.BadRequest(ctx, "User IDs tidak boleh kosong", "user_ids is required")
	}

	if err := c.service.RemoveMembers(publicID, req.UserIDs); err != nil {
		return utils.BadRequest(ctx, "Gagal menghapus member", err.Error())
	}

	return utils.Success(ctx, "Member berhasil dihapus", nil)
}

// GetMyBoardPaginated godoc
// @Summary      Get my boards
// @Description  Get all boards for the authenticated user with pagination
// @Tags         boards
// @Produce      json
// @Security     BearerAuth
// @Param        page    query     int     false  "Page number"    default(1)
// @Param        limit   query     int     false  "Items per page" default(10)
// @Param        filter  query     string  false  "Filter by title"
// @Param        sort    query     string  false  "Sort field, prefix with - for DESC (e.g. -id)"
// @Success      200     {object}  utils.ResponsePaginated{data=[]models.Board}
// @Failure      500     {object}  utils.Response
// @Router       /api/v1/boards/my [get]
func (c *BoardController) GetMyBoardPaginated(ctx fiber.Ctx) error {
	token, ok := ctx.Locals("user").(*jwt.Token)
	if !ok {
		return utils.Unauthorized(ctx, "Token tidak valid", "invalid token type")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return utils.Unauthorized(ctx, "Token tidak valid", "invalid claims type")
	}
	userID, ok := claims["pub_id"].(string)
	if !ok {
		return utils.Unauthorized(ctx, "Token tidak valid", "pub_id missing")
	}

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit
	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")

	boards, total, err := c.service.GetAllByUserPaginate(userID, filter, sort, limit, offset)
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal mengambil data board", err.Error())
	}

	meta := utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		Filter:    filter,
		Sort:      sort,
	}

	if total == 0 {
		return utils.NotFoundPagination(ctx, "Data board tidak ditemukan", boards, meta)
	}

	return utils.SuccessPagination(ctx, "Data board berhasil diambil", boards, meta)
}
