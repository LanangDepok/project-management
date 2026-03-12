package controllers

import (
	"math"
	"strconv"

	"github.com/LanangDepok/project-management/models"
	"github.com/LanangDepok/project-management/services"
	"github.com/LanangDepok/project-management/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type UserController struct {
	service services.UserService
}

func NewUserController(s services.UserService) *UserController {
	return &UserController{service: s}
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.RegisterRequest  true  "Register payload"
// @Success      200   {object}  utils.Response{data=models.UserResponse}
// @Failure      400   {object}  utils.Response
// @Router       /v1/auth/register [post]
func (c *UserController) Register(ctx fiber.Ctx) error {
	var req models.RegisterRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := c.service.Register(user); err != nil {
		return utils.BadRequest(ctx, "Registrasi Gagal", err.Error())
	}

	var userResp models.UserResponse
	_ = copier.Copy(&userResp, user)
	return utils.Success(ctx, "Register berhasil", userResp)
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.LoginRequest  true  "Login payload"
// @Success      200   {object}  utils.Response{data=models.LoginResponse}
// @Failure      401   {object}  utils.Response
// @Router       /v1/auth/login [post]
func (c *UserController) Login(ctx fiber.Ctx) error {
	var req models.LoginRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return utils.BadRequest(ctx, "Invalid Request", err.Error())
	}

	user, err := c.service.Login(req.Email, req.Password)
	if err != nil {
		return utils.Unauthorized(ctx, "Login Gagal", err.Error())
	}

	token, err := utils.GenerateToken(user.InternalID, user.Role, user.Email, user.PublicID)
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal generate token", err.Error())
	}
	refreshToken, err := utils.GenerateRefreshToken(user.InternalID)
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal generate refresh token", err.Error())
	}

	var userResp models.UserResponse
	_ = copier.Copy(&userResp, user)

	return utils.Success(ctx, "Login berhasil", models.LoginResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		User:         userResp,
	})
}

// GetUser godoc
// @Summary      Get user by ID
// @Description  Retrieve a single user by their public UUID
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User Public UUID"
// @Success      200  {object}  utils.Response{data=models.UserResponse}
// @Failure      404  {object}  utils.Response
// @Router       /api/v1/users/{id} [get]
func (c *UserController) GetUser(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := c.service.GetByPublicID(id)
	if err != nil {
		return utils.NotFound(ctx, "Data Not Found", err.Error())
	}

	var userResp models.UserResponse
	if err = copier.Copy(&userResp, user); err != nil {
		return utils.InternalServerError(ctx, "Error parsing data", err.Error())
	}
	return utils.Success(ctx, "Data berhasil ditemukan", userResp)
}

// GetUserPagination godoc
// @Summary      Get paginated users
// @Description  Retrieve users with pagination, filtering, and sorting
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        page    query     int     false  "Page number"    default(1)
// @Param        limit   query     int     false  "Items per page" default(10)
// @Param        filter  query     string  false  "Filter by name or email"
// @Param        sort    query     string  false  "Sort field, prefix with - for DESC (e.g. -id)"
// @Success      200     {object}  utils.ResponsePaginated{data=[]models.UserResponse}
// @Failure      404     {object}  utils.ResponsePaginated
// @Router       /api/v1/users/page [get]
func (c *UserController) GetUserPagination(ctx fiber.Ctx) error {
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

	users, total, err := c.service.GetAllPagination(filter, sort, limit, offset)
	if err != nil {
		return utils.BadRequest(ctx, "Gagal Mengambil Data", err.Error())
	}

	var userResp []models.UserResponse
	_ = copier.Copy(&userResp, &users)

	meta := utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		Total:     int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(limit))),
		Filter:    filter,
		Sort:      sort,
	}

	if total == 0 {
		return utils.NotFoundPagination(ctx, "Data pengguna tidak ditemukan", userResp, meta)
	}
	return utils.SuccessPagination(ctx, "Data ditemukan", userResp, meta)
}

// UpdateUser godoc
// @Summary      Update user
// @Description  Update user data by public UUID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      string                  true  "User Public UUID"
// @Param        body  body      models.UpdateUserRequest true  "Update payload"
// @Success      200   {object}  utils.Response{data=models.UserResponse}
// @Failure      400   {object}  utils.Response
// @Router       /api/v1/users/{id} [put]
func (c *UserController) UpdateUser(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	publicID, err := uuid.Parse(id)
	if err != nil {
		return utils.BadRequest(ctx, "Format ID tidak valid", err.Error())
	}

	var req models.UpdateUserRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}

	user := &models.User{
		PublicID: publicID,
		Name:     req.Name,
	}

	if err := c.service.Update(user); err != nil {
		return utils.BadRequest(ctx, "Gagal Update Data", err.Error())
	}

	userUpdated, err := c.service.GetByPublicID(id)
	if err != nil {
		return utils.InternalServerError(ctx, "Gagal Ambil Data", err.Error())
	}

	var userResp models.UserResponse
	if err = copier.Copy(&userResp, userUpdated); err != nil {
		return utils.InternalServerError(ctx, "Error parsing data", err.Error())
	}
	return utils.Success(ctx, "Berhasil update data", userResp)
}

// DeleteUser godoc
// @Summary      Delete user
// @Description  Soft-delete a user by internal ID
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "User Internal ID"
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /api/v1/users/{id} [delete]
func (c *UserController) DeleteUser(ctx fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return utils.BadRequest(ctx, "Format ID tidak valid", err.Error())
	}
	if err := c.service.Delete(uint(id)); err != nil {
		return utils.InternalServerError(ctx, "Gagal Menghapus Data", err.Error())
	}
	return utils.Success(ctx, "Berhasil menghapus data", fiber.Map{"id": id})
}
