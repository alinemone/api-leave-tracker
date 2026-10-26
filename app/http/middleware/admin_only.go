package middleware

import (
	"github.com/goravel/framework/contracts/http"

	"leave/app/helpers"
)

// AdminOnly must run after WebAuth. It allows only admin users through;
// everyone else is redirected to their leaves page.
func AdminOnly() http.Middleware {
	return func(ctx http.Context) {
		user, err := helpers.CurrentUser(ctx)
		if err != nil {
			ctx.Response().Redirect(http.StatusFound, "/login").Abort()
			return
		}
		if !user.IsAdmin {
			ctx.Response().Redirect(http.StatusFound, "/leaves-me?error=forbidden").Abort()
			return
		}
		ctx.Request().Next()
	}
}
