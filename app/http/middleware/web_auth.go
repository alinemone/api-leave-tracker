package middleware

import (
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"leave/app/models"
)

// WebAuth protects server-rendered pages. It loads the logged-in user from the
// session and stores it on the context (key "user") for helpers.CurrentUser.
// Unauthenticated visitors are redirected to the login page.
func WebAuth() http.Middleware {
	return func(ctx http.Context) {
		session := ctx.Request().Session()
		if session == nil {
			ctx.Response().Redirect(http.StatusFound, "/login").Abort()
			return
		}

		raw := session.Get("user_id")
		idStr, ok := raw.(string)
		if !ok || idStr == "" {
			ctx.Response().Redirect(http.StatusFound, "/login").Abort()
			return
		}

		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			session.Forget("user_id")
			ctx.Response().Redirect(http.StatusFound, "/login").Abort()
			return
		}

		var user models.User
		if err := facades.Orm().Query().Where("id", uint(id)).First(&user); err != nil || user.ID == 0 {
			session.Forget("user_id")
			ctx.Response().Redirect(http.StatusFound, "/login").Abort()
			return
		}

		ctx.WithValue("user", &user)
		ctx.Request().Next()
	}
}
