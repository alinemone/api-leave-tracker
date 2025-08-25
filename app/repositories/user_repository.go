package repositories

import (
	"github.com/goravel/framework/facades"
	"leave/app/interfaces"
	"leave/app/models"
)

type userRepository struct{}

func NewUserRepository() interfaces.UserInterface {
	return &userRepository{}
}

func (r *userRepository) All() ([]models.User, error) {
	var users []models.User
	err := facades.Orm().Query().Find(&users)
	return users, err
}

func (r *userRepository) List(page int, perPage int) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	err := facades.Orm().Query().Paginate(page, perPage, &users, &total)
	return users, total, err
}
