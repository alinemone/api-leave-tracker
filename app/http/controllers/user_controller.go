package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"leave/app/http/responses"
	response "leave/app/http/responses/user"
	"leave/app/interfaces"
	"leave/app/repositories"
)

type UserController struct {
	repository interfaces.UserInterface
}

func NewUserController() *UserController {
	return &UserController{
		repository: repositories.NewUserRepository(),
	}
}

func (r *UserController) All(ctx http.Context) http.Response {
	allUsers, _ := r.repository.All()

	users := response.NewLiteUserCollection(allUsers)

	return ctx.Response().Success().Json(users)
}

func (r *UserController) List(ctx http.Context) http.Response {
	page := ctx.Request().QueryInt("page", 1)
	perPage := ctx.Request().QueryInt("per_page", 10)
	data, total, _ := r.repository.List(page, perPage)

	users := response.NewLiteUserCollection(data)
	res := responses.NewPagination(users, page, perPage, total)

	return ctx.Response().Success().Json(res)
}
