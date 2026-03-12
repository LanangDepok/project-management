package repositories

import (
	"strings"

	"github.com/LanangDepok/project-management/config"
	"github.com/LanangDepok/project-management/models"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	FindByPublicID(publicID string) (*models.User, error)
	FindAllPagination(filter, sort string, limit, offset int) ([]models.User, int64, error)
	Update(user *models.User) error
	Delete(id uint) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByPublicID(publicID string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("public_id = ?", publicID).First(&user).Error
	return &user, err
}

func (r *userRepository) FindAllPagination(filter, sort string, limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	db := config.DB.Model(&models.User{})

	if filter != "" {
		pattern := "%" + filter + "%"
		db = db.Where("name ILIKE ? OR email ILIKE ?", pattern, pattern)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sort != "" {
		// Normalize sort=id / sort=-id → internal column name
		col := sort
		direction := "ASC"
		if strings.HasPrefix(col, "-") {
			col = strings.TrimPrefix(col, "-")
			direction = "DESC"
		}
		if col == "id" {
			col = "internal_id"
		}
		db = db.Order(col + " " + direction)
	}

	err := db.Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *userRepository) Update(user *models.User) error {
	return config.DB.Model(&models.User{}).
		Where("public_id = ?", user.PublicID).
		Updates(map[string]interface{}{"name": user.Name}).Error
}

func (r *userRepository) Delete(id uint) error {
	return config.DB.Delete(&models.User{}, id).Error
}
