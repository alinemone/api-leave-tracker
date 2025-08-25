package controllers

import (
	"leave/app/helpers"
	"leave/app/http/requests/auth"
	"leave/app/http/responses"
	response "leave/app/http/responses/user"
	"leave/app/models"
	"leave/services/gap"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	// Dependent services
}

func NewAuthController() *AuthController {
	return &AuthController{
		// Inject services
	}
}

func (r *AuthController) Login(ctx http.Context) http.Response {
	var req auth.LoginRequest

	errors, _ := ctx.Request().ValidateRequest(&req)

	if errors != nil {
		return ctx.Response().Json(
			http.StatusUnprocessableEntity,
			responses.NewErrorResponse(errors),
		)
	}

	var user models.User

	if err := facades.Orm().Query().Where("mobile", req.Mobile).First(&user); err != nil {
		return ctx.Response().Json(
			http.StatusInternalServerError,
			responses.NewErrorResponse(err),
		)
	}

	if user.ID == 0 {
		return ctx.Response().Json(
			http.StatusInternalServerError,
			responses.NewErrorResponse("User not found"),
		)
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return ctx.Response().Json(
			http.StatusUnauthorized,
			responses.NewErrorResponse(err),
		)
	}

	var (
		token string
	)

	token, err = facades.Auth(ctx).LoginUsingID(user.ID)

	if err != nil {
		return ctx.Response().Json(
			http.StatusInternalServerError,
			responses.NewErrorResponse(err),
		)
	}

	return ctx.Response().Success().Json(
		response.NewLoginResource(user, token),
	)
}

func (r *AuthController) Register(ctx http.Context) http.Response {

	var req auth.RegisterRequest

	errors, err := ctx.Request().ValidateRequest(&req)

	if err != nil {
		return ctx.Response().Json(
			http.StatusInternalServerError,
			responses.NewErrorResponse(err),
		)
	}

	if errors != nil {
		return ctx.Response().Json(
			http.StatusUnprocessableEntity,
			responses.NewErrorResponse(errors),
		)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return ctx.Response().Json(
			http.StatusInternalServerError,
			responses.NewErrorResponse(err),
		)

	}

	user := models.User{
		Name:     req.Name,
		IsActive: req.IsActive,
		Mobile:   req.Mobile,
		Password: string(hashedPassword),
	}

	err = facades.Orm().Query().Create(&user)

	if err != nil {
		return ctx.Response().Json(
			http.StatusInternalServerError,
			responses.NewErrorResponse(err),
		)
	}

	return ctx.Response().Success().Json(
		response.NewLiteUserResource(user),
	)
}

func (r *AuthController) RequestCode(ctx http.Context) http.Response {
	var req auth.ResetCodeRequest

	errors, _ := ctx.Request().ValidateRequest(&req)

	if errors != nil {
		return ctx.Response().Json(
			http.StatusUnprocessableEntity,
			responses.NewErrorResponse(errors),
		)
	}

	var user models.User
	if err := facades.Orm().Query().Where("mobile", req.Mobile).First(&user); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": err.Error(),
		})
	}

	if user.ID == 0 {
		return ctx.Response().Json(
			http.StatusInternalServerError,
			responses.NewErrorResponse("User not found"),
		)
	}

	code := helpers.RandomNumbers(6)

	gapClient := gap.NewClient()
	err := gapClient.SendDirectMessage(req.GapUserName, code)

	if err != nil {
		return ctx.Response().Json(
			http.StatusInternalServerError,
			responses.NewErrorResponse(err),
		)
	}

	cacheKey := "forgot_password_code:" + req.GapUserName
	err = facades.Cache().Put(cacheKey, code, 5*time.Minute)

	if err != nil {
		return nil
	}

	return ctx.Response().Success().Json(
		responses.NewPublicResponse(true),
	)
}

func (r *AuthController) ForgetPassword(ctx http.Context) http.Response {
	var req auth.ForgetPasswordRequest

	errors, _ := ctx.Request().ValidateRequest(&req)

	if errors != nil {
		return ctx.Response().Json(
			http.StatusUnprocessableEntity,
			responses.NewErrorResponse(errors),
		)
	}

	var user models.User
	if err := facades.Orm().Query().Where("mobile", req.Mobile).First(&user); err != nil {
		return ctx.Response().Json(
			http.StatusInternalServerError,
			responses.NewErrorResponse(err),
		)
	}

	if user.ID == 0 {
		return ctx.Response().Json(
			http.StatusInternalServerError,
			responses.NewErrorResponse("User not found"),
		)
	}

	value := facades.Cache().Get("forgot_password_code:" + req.GapUserName)

	if value == nil {
		return ctx.Response().Json(
			http.StatusInternalServerError,
			responses.NewErrorResponse("code expired try again reset code"),
		)
	}

	if value != req.Code {
		return ctx.Response().Json(
			http.StatusUnprocessableEntity,
			responses.NewErrorResponse("code not match"),
		)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return ctx.Response().Json(
			http.StatusInternalServerError,
			responses.NewErrorResponse(err),
		)
	}

	user.Password = string(hashedPassword)
	err = facades.Orm().Query().Save(&user)
	if err != nil {
		return ctx.Response().Json(
			http.StatusInternalServerError,
			responses.NewErrorResponse(err),
		)
	}

	return ctx.Response().Success().Json(
		responses.NewPublicResponse(true),
	)
}

func (r *AuthController) Whoami(ctx http.Context) http.Response {
	currentUser, _ := helpers.CurrentUser(ctx)

	return ctx.Response().Success().Json(
		response.NewUserResource(*currentUser),
	)
}
