package middleware

import (
	"github.com/goravel/framework/auth"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/facades"
	"strings"
)

func Auth1() http.Middleware {
	return func(ctx http.Context) {
		guard := facades.Config().GetString("auth.defaults.guard")

		if ctx.Request().Header("Guard") != "" {
			guard = ctx.Request().Header("Guard")
		}

		token := ctx.Request().Header("Authorization")

		if token == "" {
			ctx.Request().AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if _, err := facades.Auth(ctx).Guard(guard).Parse(token); err != nil {
			if errors.Is(err, auth.ErrorTokenExpired) {
				token, err := facades.Auth(ctx).Guard(guard).Refresh()
				if err != nil {
					ctx.Request().AbortWithStatus(http.StatusUnauthorized)
					return
				}

				token = "Bearer " + token

			} else {
				ctx.Request().AbortWithStatus(http.StatusUnauthorized)
				return
			}

		}

		ctx.Response().Header("Authorization", token)
		ctx.Request().Next()
	}
}

func Auth() http.Middleware {
	return func(ctx http.Context) {
		guard := facades.Config().GetString("auth.defaults.guard")

		authHeader := ctx.Request().Header("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.Request().AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		user, err := facades.Auth(ctx).Guard(guard).Parse(token)
		if err != nil {
			if err == auth.ErrorTokenExpired {
				ctx.Request().AbortWithStatus(http.StatusUnauthorized)
				return
			}
			ctx.Request().AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.WithValue("user", user)
		ctx.Request().Next()
	}
}
