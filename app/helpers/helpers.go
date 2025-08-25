package helpers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"golang.org/x/exp/rand"
	"leave/app/models"
	"time"
)

func CurrentUser(ctx http.Context) (*models.User, error) {
	var user models.User
	err := facades.Auth(ctx).User(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func RandomNumbers(length int) string {
	rand.Seed(uint64(time.Now().UnixNano()))
	digits := "0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = digits[rand.Intn(len(digits))]
	}
	return string(result)
}
