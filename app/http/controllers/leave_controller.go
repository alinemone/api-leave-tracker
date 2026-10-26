package controllers

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"leave/app/helpers"
	leaveRequest "leave/app/http/requests/leave"
	"leave/app/interfaces"
	"leave/app/repositories"
	"leave/services/notify"
)

const perPage = 10

type LeaveController struct {
	repository interfaces.LeaveInterface
	teams      *repositories.TeamRepository
}

func NewLeaveController() *LeaveController {
	return &LeaveController{
		repository: repositories.NewLeaveRepository(),
		teams:      repositories.NewTeamRepository(),
	}
}

// MyLeaves renders the current user's leaves, their own totals, and the create form.
func (r *LeaveController) MyLeaves(ctx http.Context) http.Response {
	user, _ := helpers.CurrentUser(ctx)
	page := ctx.Request().QueryInt("page", 1)

	leaves, total, _ := r.repository.MyLeaves(user.ID, page, perPage)
	teams, _ := r.teams.MemberTeams(user.ID)

	jYear, jMonth := helpers.CurrentJalali()
	yf, yt := helpers.JalaliYearRange(jYear)
	mf, mt := helpers.JalaliMonthRange(jYear, jMonth)
	yearTot, _ := r.repository.MyTotals(user.ID, yf, yt)
	monthTot, _ := r.repository.MyTotals(user.ID, mf, mt)

	return ctx.Response().View().Make("leaves_me.tmpl", map[string]any{
		"title":      facades.Lang(ctx).Get("nav.my_leaves"),
		"active":     "me",
		"user":       user,
		"leaves":     leaves,
		"teams":      teams,
		"pagination": newPagination(page, total),
		"jYear":      jYear,
		"jMonthName": helpers.JalaliMonthName(jMonth),
		"monthDays":  monthTot.TotalDays,
		"monthHours": helpers.FormatHours(monthTot.TotalHourSeconds),
		"yearDays":   yearTot.TotalDays,
		"yearHours":  helpers.FormatHours(yearTot.TotalHourSeconds),
		"error":      ctx.Request().Query("error"),
		"success":    ctx.Request().Query("success"),
	})
}

// Store validates and creates a personal or team leave.
func (r *LeaveController) Store(ctx http.Context) http.Response {
	var req leaveRequest.CreateRequest

	errors, err := ctx.Request().ValidateRequest(&req)
	if err != nil {
		return redirectWithError(ctx, facades.Lang(ctx).Get("leave.errors.validate_failed"))
	}
	if errors != nil {
		return redirectWithError(ctx, errors.One())
	}

	user, err := helpers.CurrentUser(ctx)
	if err != nil {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	var teamID *uint
	if ctx.Request().Input("target") == "team" {
		tid := uint(ctx.Request().InputInt("team_id"))
		if tid == 0 || !r.teams.IsMember(user.ID, tid) {
			return redirectWithError(ctx, facades.Lang(ctx).Get("leave.errors.not_member"))
		}
		teamID = &tid
	}

	leave, err := r.repository.Create(user.ID, req.Type, req.StartAt, req.EndAt, teamID)
	if err != nil || leave.ID == 0 {
		return redirectWithError(ctx, facades.Lang(ctx).Get("leave.errors.create"))
	}

	// Notify the team's webhook (if configured) about the new team leave.
	if teamID != nil {
		if team, e := r.teams.Find(*teamID); e == nil && team.WebhookURL != "" {
			notify.SendAsync(team.WebhookURL, leaveNotificationText(user.Name, team.Name, leave.Type, leave.StartAt, leave.EndAt))
		}
	}

	return ctx.Response().Redirect(http.StatusFound, "/leaves-me?success=1")
}

// Delete soft-deletes a leave owned by the current user.
func (r *LeaveController) Delete(ctx http.Context) http.Response {
	currentUser, err := helpers.CurrentUser(ctx)
	if err != nil {
		return ctx.Response().Redirect(http.StatusFound, "/login")
	}

	id64, err := strconv.ParseUint(ctx.Request().Input("id"), 10, 64)
	if err != nil || id64 == 0 {
		return redirectWithError(ctx, facades.Lang(ctx).Get("leave.errors.invalid_id"))
	}

	// Read the leave first so we can notify its team after deletion.
	leave, _ := r.repository.FindByUser(currentUser.ID, uint(id64))

	if err := r.repository.Delete(currentUser.ID, uint(id64)); err != nil {
		return redirectWithError(ctx, facades.Lang(ctx).Get("leave.errors.delete"))
	}

	// Notify the team's webhook (if configured) that a team leave was cancelled.
	if leave.ID != 0 && leave.TeamID != nil {
		if team, e := r.teams.Find(*leave.TeamID); e == nil && team.WebhookURL != "" {
			notify.SendAsync(team.WebhookURL, leaveCancelText(currentUser.Name, team.Name, leave.Type, leave.StartAt, leave.EndAt))
		}
	}

	return ctx.Response().Redirect(http.StatusFound, "/leaves-me?success=deleted")
}

// pagination holds the minimal data the templates need to render page links.
type pagination struct {
	Page     int
	PerPage  int
	Total    int64
	HasPrev  bool
	HasNext  bool
	PrevPage int
	NextPage int
}

func newPagination(page int, total int64) pagination {
	if page < 1 {
		page = 1
	}
	return pagination{
		Page:     page,
		PerPage:  perPage,
		Total:    total,
		HasPrev:  page > 1,
		HasNext:  int64(page*perPage) < total,
		PrevPage: page - 1,
		NextPage: page + 1,
	}
}

func redirectWithError(ctx http.Context, message string) http.Response {
	return ctx.Response().Redirect(http.StatusFound, "/leaves-me?error="+url.QueryEscape(message))
}

// leaveNotificationText builds the team webhook message for a new leave.
func leaveNotificationText(userName, teamName, leaveType string, start, end time.Time) string {
	if leaveType == "hourly" {
		return fmt.Sprintf("⏰ %s برای تیم «%s» مرخصی ساعتی در تاریخ %s از %s تا %s ثبت کرد.",
			userName, teamName, helpers.JalaliLong(start), helpers.Clock(start), helpers.Clock(end))
	}
	return fmt.Sprintf("📅 %s برای تیم «%s» مرخصی روزانه در تاریخ %s ثبت کرد.",
		userName, teamName, helpers.JalaliLong(start))
}

// leaveCancelText builds the team webhook message for a cancelled leave.
func leaveCancelText(userName, teamName, leaveType string, start, end time.Time) string {
	if leaveType == "hourly" {
		return fmt.Sprintf("❌ %s مرخصی ساعتی خود در تیم «%s» را در تاریخ %s از %s تا %s لغو کرد.",
			userName, teamName, helpers.JalaliLong(start), helpers.Clock(start), helpers.Clock(end))
	}
	return fmt.Sprintf("❌ %s مرخصی روزانه خود در تیم «%s» را در تاریخ %s لغو کرد.",
		userName, teamName, helpers.JalaliLong(start))
}
