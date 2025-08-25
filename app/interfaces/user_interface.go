package interfaces

import (
	"leave/app/models"
)

type UserInterface interface {
	All() ([]models.User, error)
	List(page int, perPage int) ([]models.User, int64, error)
}
