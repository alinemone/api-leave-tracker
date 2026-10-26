package models

import (
	"github.com/goravel/framework/database/orm"
)

type User struct {
	orm.Model
	Name string
	// BasalamUserId is the unique Basalam account id this user is linked to.
	// Pointer so manually-created users store NULL (NULLs don't collide in the
	// unique index) instead of all sharing 0.
	BasalamUserId *uint `gorm:"uniqueIndex"`
	IsActive      bool
	// IsAdmin grants company management and user administration rights.
	IsAdmin bool
	Mobile  string `gorm:"uniqueIndex"`
	// Password is unused now that authentication is handled via Basalam OAuth.
	Password string
}
