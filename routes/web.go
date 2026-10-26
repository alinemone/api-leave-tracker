package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/path"

	"leave/app/http/controllers"
	"leave/app/http/middleware"
)

func Web() {
	authController := controllers.NewAuthController()
	leaveController := controllers.NewLeaveController()
	teamController := controllers.NewTeamController()
	adminController := controllers.NewAdminController()
	dashboardController := controllers.NewDashboardController()
	profileController := controllers.NewProfileController()

	// Static assets (fonts + logo + vendored js/css placed in public/).
	facades.Route().Static("/fonts", path.Public("fonts"))
	facades.Route().Static("/js", path.Public("js"))
	facades.Route().Static("/css", path.Public("css"))
	facades.Route().StaticFile("/logo.png", path.Public("logo.png"))

	// Landing: send visitors to their leaves (WebAuth bounces guests to /login).
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		return ctx.Response().Redirect(http.StatusFound, "/leaves-me")
	})

	// Public auth + manual registration.
	facades.Route().Get("/login", authController.ShowLogin)
	facades.Route().Post("/login", authController.Login)
	facades.Route().Get("/register", authController.ShowRegister)
	facades.Route().Post("/register", authController.Register)
	facades.Route().Get("/auth/basalam/redirect", authController.BasalamRedirect)
	facades.Route().Get("/auth/basalam/callback", authController.BasalamCallback)
	facades.Route().Post("/logout", authController.Logout)

	// Authenticated pages (any logged-in user).
	facades.Route().Middleware(middleware.WebAuth()).Group(func(router route.Router) {
		router.Get("/leaves-me", leaveController.MyLeaves)
		router.Post("/leaves", leaveController.Store)
		router.Post("/leaves/delete", leaveController.Delete)

		router.Get("/profile", profileController.Show)
		router.Post("/profile", profileController.Update)

		router.Get("/dashboard", dashboardController.Index)
		router.Post("/teams/{id}/reset", teamController.Reset)

		router.Get("/teams", teamController.Index)
		router.Post("/teams", teamController.Store)
		router.Get("/teams/{id}", teamController.Show)
		router.Get("/teams/{id}/edit", teamController.Edit)
		router.Post("/teams/{id}/update", teamController.Update)
		router.Post("/teams/{id}/delete", teamController.Delete)
		router.Post("/teams/{id}/join", teamController.RequestJoin)
		router.Post("/teams/{id}/cancel-join", teamController.CancelJoin)
		router.Post("/teams/{id}/members", teamController.AddMember)
		router.Post("/teams/{id}/members/approve", teamController.ApproveMember)
		router.Post("/teams/{id}/members/set-manager", teamController.SetManager)
		router.Post("/teams/{id}/members/delete", teamController.RemoveMember)
	})

	// Admin-only (user administration + annual reset).
	facades.Route().Middleware(middleware.WebAuth(), middleware.AdminOnly()).Group(func(router route.Router) {
		router.Get("/admin/users", adminController.Users)
		router.Post("/admin/users/{id}/toggle-admin", adminController.ToggleAdmin)
		router.Post("/admin/reset", adminController.Reset)
	})
}
