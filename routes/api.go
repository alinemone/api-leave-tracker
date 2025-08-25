package routes

import (
	"leave/app/http/middleware"

	"github.com/goravel/framework/facades"

	"leave/app/http/controllers"
)

func Api() {
	router := facades.Route().Prefix("api/v1")

	userController := controllers.NewUserController()
	router.Get("/users", userController.All)
	router.Get("/users/list", userController.List)

	authController := controllers.NewAuthController()
	router.Post("auth/login", authController.Login)
	router.Post("auth/register", authController.Register)
	router.Post("auth/request-code", authController.RequestCode)
	router.Post("auth/forget-password", authController.ForgetPassword)
	router.Middleware(middleware.Auth()).Get("auth/whoami", authController.Whoami)

	leaveController := controllers.NewLeaveController()
	router.Middleware(middleware.Auth()).Get("leaves", leaveController.Report)
	router.Middleware(middleware.Auth()).Get("leaves/me", leaveController.ListUserLeaves)
	router.Middleware(middleware.Auth()).Post("leaves", leaveController.Store)
	router.Middleware(middleware.Auth()).Delete("leaves", leaveController.Delete)
}
