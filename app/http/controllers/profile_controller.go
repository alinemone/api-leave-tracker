package controllers

import (
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"golang.org/x/crypto/bcrypt"

	"leave/app/helpers"
	"leave/app/models"
)

type ProfileController struct {
}

func NewProfileController() *ProfileController {
	return &ProfileController{}
}

var profileErrors = map[string]bool{
	"required":       true,
	"mobile_taken":   true,
	"password_short": true,
	"failed":         true,
}

// Show renders the current user's editable profile.
func (r *ProfileController) Show(ctx http.Context) http.Response {
	user, _ := helpers.CurrentUser(ctx)

	var errMessage string
	if code := ctx.Request().Query("error"); profileErrors[code] {
		errMessage = facades.Lang(ctx).Get("profile.errors." + code)
	}

	return ctx.Response().View().Make("profile.tmpl", map[string]any{
		"title":   facades.Lang(ctx).Get("profile.title"),
		"active":  "profile",
		"user":    user,
		"error":   errMessage,
		"success": ctx.Request().Query("success"),
	})
}

// Update saves the user's own name, mobile and (optionally) a new password.
func (r *ProfileController) Update(ctx http.Context) http.Response {
	user, err := helpers.CurrentUser(ctx)
	if err != nil {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	name := strings.TrimSpace(ctx.Request().Input("name"))
	mobile := helpers.EnglishDigits(strings.TrimSpace(ctx.Request().Input("mobile")))
	password := ctx.Request().Input("password")

	if name == "" || mobile == "" {
		return ctx.Response().Redirect(http.StatusFound, "/profile?error=required")
	}

	// Mobile must stay unique (ignoring this user's own row).
	var other models.User
	_ = facades.Orm().Query().Where("mobile", mobile).Where("id <> ?", user.ID).First(&other)
	if other.ID != 0 {
		return ctx.Response().Redirect(http.StatusFound, "/profile?error=mobile_taken")
	}

	user.Name = name
	user.Mobile = mobile

	if password != "" {
		if len(password) < 4 {
			return ctx.Response().Redirect(http.StatusFound, "/profile?error=password_short")
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return ctx.Response().Redirect(http.StatusFound, "/profile?error=failed")
		}
		user.Password = string(hashed)
	}

	if err := facades.Orm().Query().Save(user); err != nil {
		facades.Log().Error(err)
		return ctx.Response().Redirect(http.StatusFound, "/profile?error=failed")
	}

	return ctx.Response().Redirect(http.StatusFound, "/profile?success=1")
}
