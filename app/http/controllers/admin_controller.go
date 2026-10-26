package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"leave/app/helpers"
	"leave/app/interfaces"
	"leave/app/models"
	"leave/app/repositories"
)

type AdminController struct {
	leaves interfaces.LeaveInterface
}

func NewAdminController() *AdminController {
	return &AdminController{
		leaves: repositories.NewLeaveRepository(),
	}
}

// Users lists all users so an admin can grant/revoke admin rights.
func (r *AdminController) Users(ctx http.Context) http.Response {
	current, _ := helpers.CurrentUser(ctx)

	var users []models.User
	_ = facades.Orm().Query().OrderByDesc("is_admin").OrderBy("name").Find(&users)

	return ctx.Response().View().Make("admin_users.tmpl", map[string]any{
		"title":   facades.Lang(ctx).Get("admin.users_title"),
		"active":  "admin",
		"user":    current,
		"users":   users,
		"years":   helpers.JalaliYears(),
		"error":   ctx.Request().Query("error"),
		"success": ctx.Request().Query("success"),
	})
}

// Reset soft-deletes every leave for a chosen Jalali year (admin-wide).
func (r *AdminController) Reset(ctx http.Context) http.Response {
	year := ctx.Request().InputInt("year")
	if year == 0 {
		year, _ = helpers.CurrentJalali()
	}
	from, to := helpers.JalaliYearRange(year)
	if _, err := r.leaves.ResetAll(from, to); err != nil {
		return ctx.Response().Redirect(http.StatusFound, "/admin/users?error=failed")
	}
	return ctx.Response().Redirect(http.StatusFound, "/admin/users?success=1")
}

// ToggleAdmin flips a user's admin flag (refusing to remove the last admin).
func (r *AdminController) ToggleAdmin(ctx http.Context) http.Response {
	id := ctx.Request().RouteInt("id")

	var target models.User
	if err := facades.Orm().Query().Where("id", id).First(&target); err != nil || target.ID == 0 {
		return ctx.Response().Redirect(http.StatusFound, "/admin/users?error=notfound")
	}

	if target.IsAdmin {
		adminCount, _ := facades.Orm().Query().Model(&models.User{}).Where("is_admin", true).Count()
		if adminCount <= 1 {
			return ctx.Response().Redirect(http.StatusFound, "/admin/users?error=last_admin")
		}
	}

	if _, err := facades.Orm().Query().Model(&models.User{}).Where("id", target.ID).Update("is_admin", !target.IsAdmin); err != nil {
		return ctx.Response().Redirect(http.StatusFound, "/admin/users?error=failed")
	}

	return ctx.Response().Redirect(http.StatusFound, "/admin/users?success=1")
}
