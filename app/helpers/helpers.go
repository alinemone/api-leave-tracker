package helpers

import (
	"errors"

	"github.com/goravel/framework/contracts/http"

	"leave/app/models"
)

// CurrentUser returns the authenticated user that the WebAuth middleware stored
// on the request context.
func CurrentUser(ctx http.Context) (*models.User, error) {
	value := ctx.Value("user")
	if value == nil {
		return nil, errors.New("user not found in context")
	}

	user, ok := value.(*models.User)
	if !ok {
		return nil, errors.New("invalid user in context")
	}

	return user, nil
}
