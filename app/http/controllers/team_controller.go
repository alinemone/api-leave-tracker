package controllers

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"leave/app/helpers"
	teamresp "leave/app/http/responses/team"
	"leave/app/interfaces"
	"leave/app/models"
	"leave/app/repositories"
	"leave/services/notify"
)

type TeamController struct {
	repository *repositories.TeamRepository
	leaves     interfaces.LeaveInterface
}

func NewTeamController() *TeamController {
	return &TeamController{
		repository: repositories.NewTeamRepository(),
		leaves:     repositories.NewLeaveRepository(),
	}
}

// Index lists the user's teams and offers a create form (any user, any company).
func (r *TeamController) Index(ctx http.Context) http.Response {
	user, _ := helpers.CurrentUser(ctx)
	teams, _ := r.repository.MemberTeams(user.ID)
	joinable, _ := r.repository.JoinableTeams(user.ID)

	return ctx.Response().View().Make("teams.tmpl", map[string]any{
		"title":    facades.Lang(ctx).Get("nav.teams"),
		"active":   "teams",
		"user":     user,
		"teams":    teams,
		"joinable": joinable,
		"error":    ctx.Request().Query("error"),
		"success":  ctx.Request().Query("success"),
	})
}

// Store creates a team under a company (creator becomes manager).
func (r *TeamController) Store(ctx http.Context) http.Response {
	user, _ := helpers.CurrentUser(ctx)

	name := strings.TrimSpace(ctx.Request().Input("name"))
	webhookURL := strings.TrimSpace(ctx.Request().Input("webhook_url"))
	if name == "" {
		return teamsRedirectError(ctx, facades.Lang(ctx).Get("team.errors.name_required"))
	}

	if _, err := r.repository.Create(name, user.ID, webhookURL); err != nil {
		facades.Log().Error(err)
		return teamsRedirectError(ctx, facades.Lang(ctx).Get("team.errors.failed"))
	}

	return ctx.Response().Redirect(http.StatusFound, "/teams?success=1")
}

// Show is the team detail page (members + add member + leaves report). Managers only.
func (r *TeamController) Show(ctx http.Context) http.Response {
	user, _ := helpers.CurrentUser(ctx)
	teamID := uint(ctx.Request().RouteInt("id"))

	if !r.repository.IsManager(user.ID, teamID) {
		return ctx.Response().Redirect(http.StatusFound, "/teams?error="+url.QueryEscape(facades.Lang(ctx).Get("team.errors.forbidden")))
	}

	team, err := r.repository.Find(teamID)
	if err != nil || team.ID == 0 {
		return ctx.Response().Redirect(http.StatusFound, "/teams?error="+url.QueryEscape(facades.Lang(ctx).Get("team.errors.notfound")))
	}

	members, _ := r.repository.Members(teamID)
	pending, _ := r.repository.PendingRequests(teamID)
	page := ctx.Request().QueryInt("page", 1)
	leaves, total, _ := r.leaves.TeamLeaves(teamID, page, perPage)

	q := strings.TrimSpace(ctx.Request().Query("q"))
	var searchResults []teamresp.MemberRow
	if q != "" {
		searchResults, _ = r.repository.SearchAddableUsers(teamID, q)
	}

	return ctx.Response().View().Make("team_show.tmpl", map[string]any{
		"title":         team.Name,
		"active":        "teams",
		"user":          user,
		"team":          team,
		"members":       members,
		"pending":       pending,
		"q":             q,
		"searchResults": searchResults,
		"leaves":        leaves,
		"pagination":    newPagination(page, total),
		"error":         ctx.Request().Query("error"),
		"success":       ctx.Request().Query("success"),
	})
}

// Edit renders the team rename form (managers only).
func (r *TeamController) Edit(ctx http.Context) http.Response {
	user, _ := helpers.CurrentUser(ctx)
	teamID := uint(ctx.Request().RouteInt("id"))

	if !r.repository.IsManager(user.ID, teamID) {
		return ctx.Response().Redirect(http.StatusFound, "/teams?error="+url.QueryEscape(facades.Lang(ctx).Get("team.errors.forbidden")))
	}

	team, err := r.repository.Find(teamID)
	if err != nil || team.ID == 0 {
		return ctx.Response().Redirect(http.StatusFound, "/teams?error="+url.QueryEscape(facades.Lang(ctx).Get("team.errors.notfound")))
	}

	return ctx.Response().View().Make("team_edit.tmpl", map[string]any{
		"title":  facades.Lang(ctx).Get("team.edit_title"),
		"active": "teams",
		"user":   user,
		"team":   team,
	})
}

// Update saves the team's new name (managers only).
func (r *TeamController) Update(ctx http.Context) http.Response {
	user, _ := helpers.CurrentUser(ctx)
	teamID := uint(ctx.Request().RouteInt("id"))

	if !r.repository.IsManager(user.ID, teamID) {
		return teamsRedirectError(ctx, facades.Lang(ctx).Get("team.errors.forbidden"))
	}

	name := strings.TrimSpace(ctx.Request().Input("name"))
	webhookURL := strings.TrimSpace(ctx.Request().Input("webhook_url"))
	if name == "" {
		return teamsRedirectError(ctx, facades.Lang(ctx).Get("team.errors.name_required"))
	}

	if err := r.repository.Update(teamID, name, webhookURL); err != nil {
		return teamsRedirectError(ctx, facades.Lang(ctx).Get("team.errors.failed"))
	}

	return ctx.Response().Redirect(http.StatusFound, "/teams?success=1")
}

// Delete soft-deletes a team (managers only).
func (r *TeamController) Delete(ctx http.Context) http.Response {
	user, _ := helpers.CurrentUser(ctx)
	teamID := uint(ctx.Request().RouteInt("id"))

	if !r.repository.IsManager(user.ID, teamID) {
		return teamsRedirectError(ctx, facades.Lang(ctx).Get("team.errors.forbidden"))
	}

	if err := r.repository.Delete(teamID); err != nil {
		return teamsRedirectError(ctx, facades.Lang(ctx).Get("team.errors.failed"))
	}

	return ctx.Response().Redirect(http.StatusFound, "/teams?success=deleted")
}

// AddMember adds (or links) a user to the team by name + mobile. Managers only.
func (r *TeamController) AddMember(ctx http.Context) http.Response {
	user, _ := helpers.CurrentUser(ctx)
	teamID := uint(ctx.Request().RouteInt("id"))

	if !r.repository.IsManager(user.ID, teamID) {
		return teamsRedirectError(ctx, facades.Lang(ctx).Get("team.errors.forbidden"))
	}

	userID := uint(ctx.Request().InputInt("user_id"))
	isManager := ctx.Request().Input("is_manager") != ""

	// Only an existing (registered) user, picked from search, can be added.
	var member models.User
	if userID != 0 {
		_ = facades.Orm().Query().Where("id", userID).First(&member)
	}
	if member.ID == 0 {
		return r.teamRedirectError(ctx, teamID, facades.Lang(ctx).Get("member.errors.not_registered"))
	}

	if err := r.repository.AddMember(teamID, member.ID, isManager); err != nil {
		return r.teamRedirectError(ctx, teamID, facades.Lang(ctx).Get("member.errors.failed"))
	}

	return r.teamRedirectSuccess(ctx, teamID)
}

// RequestJoin lets any user request to join a team (pending manager approval).
func (r *TeamController) RequestJoin(ctx http.Context) http.Response {
	user, _ := helpers.CurrentUser(ctx)
	teamID := uint(ctx.Request().RouteInt("id"))

	team, err := r.repository.Find(teamID)
	if err != nil || team.ID == 0 {
		return teamsRedirectError(ctx, facades.Lang(ctx).Get("team.errors.notfound"))
	}

	if err := r.repository.RequestJoin(teamID, user.ID); err != nil {
		return teamsRedirectError(ctx, facades.Lang(ctx).Get("team.errors.failed"))
	}

	if team.WebhookURL != "" {
		base := strings.TrimRight(facades.Config().GetString("http.url"), "/")
		link := base + teamPath(teamID)
		notify.SendAsync(team.WebhookURL, fmt.Sprintf(
			"👤 %s درخواست عضویت در تیم «%s» را دارد.\nبرای بررسی و تأیید: %s",
			user.Name, team.Name, link,
		))
	}

	return ctx.Response().Redirect(http.StatusFound, "/teams?success=requested")
}

// ApproveMember approves a pending join request (managers only).
func (r *TeamController) ApproveMember(ctx http.Context) http.Response {
	user, _ := helpers.CurrentUser(ctx)
	teamID := uint(ctx.Request().RouteInt("id"))

	if !r.repository.IsManager(user.ID, teamID) {
		return teamsRedirectError(ctx, facades.Lang(ctx).Get("team.errors.forbidden"))
	}

	memberID := uint(ctx.Request().InputInt("user_id"))
	if memberID != 0 {
		_ = r.repository.ApproveMember(teamID, memberID)
	}

	return r.teamRedirectSuccess(ctx, teamID)
}

// RemoveMember removes a member from the team (managers only).
func (r *TeamController) RemoveMember(ctx http.Context) http.Response {
	user, _ := helpers.CurrentUser(ctx)
	teamID := uint(ctx.Request().RouteInt("id"))

	if !r.repository.IsManager(user.ID, teamID) {
		return teamsRedirectError(ctx, facades.Lang(ctx).Get("team.errors.forbidden"))
	}

	memberID := uint(ctx.Request().InputInt("user_id"))
	if memberID != 0 {
		_ = r.repository.RemoveMember(teamID, memberID)
	}

	return r.teamRedirectSuccess(ctx, teamID)
}

// SetManager promotes or demotes a member (managers only; keeps ≥1 manager).
func (r *TeamController) SetManager(ctx http.Context) http.Response {
	user, _ := helpers.CurrentUser(ctx)
	teamID := uint(ctx.Request().RouteInt("id"))

	if !r.repository.IsManager(user.ID, teamID) {
		return teamsRedirectError(ctx, facades.Lang(ctx).Get("team.errors.forbidden"))
	}

	memberID := uint(ctx.Request().InputInt("user_id"))
	makeManager := ctx.Request().Input("make") == "1"
	if memberID == 0 {
		return r.teamRedirectSuccess(ctx, teamID)
	}

	// Don't allow demoting the last remaining manager.
	if !makeManager && r.repository.CountManagers(teamID) <= 1 {
		return r.teamRedirectError(ctx, teamID, facades.Lang(ctx).Get("member.errors.last_manager"))
	}

	_ = r.repository.SetManager(teamID, memberID, makeManager)
	return r.teamRedirectSuccess(ctx, teamID)
}

// CancelJoin lets a user withdraw their own pending join request.
func (r *TeamController) CancelJoin(ctx http.Context) http.Response {
	user, _ := helpers.CurrentUser(ctx)
	teamID := uint(ctx.Request().RouteInt("id"))

	_ = r.repository.CancelJoinRequest(teamID, user.ID)

	return ctx.Response().Redirect(http.StatusFound, "/teams?success=1")
}

// Reset soft-deletes the team's leaves for a chosen Jalali year (managers only).
func (r *TeamController) Reset(ctx http.Context) http.Response {
	user, _ := helpers.CurrentUser(ctx)
	teamID := uint(ctx.Request().RouteInt("id"))

	if !r.repository.IsManager(user.ID, teamID) {
		return ctx.Response().Redirect(http.StatusFound, "/dashboard?error=forbidden")
	}

	year := ctx.Request().InputInt("year")
	if year == 0 {
		year, _ = helpers.CurrentJalali()
	}
	from, to := helpers.JalaliYearRange(year)
	if _, err := r.leaves.ResetTeam(teamID, from, to); err != nil {
		return ctx.Response().Redirect(http.StatusFound, "/dashboard?error=reset")
	}

	return ctx.Response().Redirect(http.StatusFound, "/dashboard?success=reset")
}

func (r *TeamController) teamRedirectError(ctx http.Context, teamID uint, message string) http.Response {
	return ctx.Response().Redirect(http.StatusFound, teamPath(teamID)+"?error="+url.QueryEscape(message))
}

func (r *TeamController) teamRedirectSuccess(ctx http.Context, teamID uint) http.Response {
	return ctx.Response().Redirect(http.StatusFound, teamPath(teamID)+"?success=1")
}

func teamPath(teamID uint) string {
	return "/teams/" + strconv.FormatUint(uint64(teamID), 10)
}

func teamsRedirectError(ctx http.Context, message string) http.Response {
	return ctx.Response().Redirect(http.StatusFound, "/teams?error="+url.QueryEscape(message))
}
