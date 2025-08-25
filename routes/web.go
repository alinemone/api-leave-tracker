package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"
	"github.com/goravel/framework/support/path"
)

func Web() {
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		return ctx.Response().View().Make("welcome.tmpl", map[string]any{
			"version": support.Version,
		})
	})

	facades.Route().Get("/docs", func(ctx http.Context) http.Response {
		return ctx.Response().View().Make("swagger.tmpl", map[string]any{
			"openapiURL": "/openapi/openapi.json",
		})
	})
	facades.Route().Get("/docs/index.html", func(ctx http.Context) http.Response {
		return ctx.Response().View().Make("swagger.tmpl", map[string]any{
			"openapiURL": "/openapi/openapi.json",
		})
	})

	facades.Route().Get("/openapi/openapi.json", func(ctx http.Context) http.Response {
		path := path.Public("openapi/openapi.json")
		return ctx.Response().File(path)
	})
}
