package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"leave/app/helpers"
	"leave/app/interfaces"
	"leave/app/repositories"
)

type DashboardController struct {
	teams  *repositories.TeamRepository
	leaves interfaces.LeaveInterface
}

func NewDashboardController() *DashboardController {
	return &DashboardController{
		teams:  repositories.NewTeamRepository(),
		leaves: repositories.NewLeaveRepository(),
	}
}

// Index shows the manager dashboard: per-member leave totals for the current
// Jalali month and year, for every team the user manages.
func (r *DashboardController) Index(ctx http.Context) http.Response {
	user, _ := helpers.CurrentUser(ctx)

	jYear, jMonth := helpers.CurrentJalali()
	yearFrom, yearTo := helpers.JalaliYearRange(jYear)
	monthFrom, monthTo := helpers.JalaliMonthRange(jYear, jMonth)

	myTeams, _ := r.teams.MemberTeams(user.ID)

	var teams []map[string]any
	for _, t := range myTeams {
		if !t.IsManager {
			continue
		}

		members, _ := r.teams.Members(t.TeamID)
		yearAgg, _ := r.leaves.LeaveTotalsByUser(t.TeamID, yearFrom, yearTo)
		monthAgg, _ := r.leaves.LeaveTotalsByUser(t.TeamID, monthFrom, monthTo)

		yearByUser := map[uint][2]float64{}
		for _, a := range yearAgg {
			yearByUser[a.UserID] = [2]float64{a.TotalDays, a.TotalHourSeconds}
		}
		monthByUser := map[uint][2]float64{}
		for _, a := range monthAgg {
			monthByUser[a.UserID] = [2]float64{a.TotalDays, a.TotalHourSeconds}
		}

		var rows []map[string]any
		for _, m := range members {
			y := yearByUser[m.UserID]
			mo := monthByUser[m.UserID]
			rows = append(rows, map[string]any{
				"name":       m.Name,
				"monthDays":  mo[0],
				"monthHours": helpers.FormatHours(mo[1]),
				"yearDays":   y[0],
				"yearHours":  helpers.FormatHours(y[1]),
			})
		}

		teams = append(teams, map[string]any{
			"teamID":  t.TeamID,
			"name":    t.Name,
			"members": rows,
		})
	}

	return ctx.Response().View().Make("dashboard.tmpl", map[string]any{
		"title":      facades.Lang(ctx).Get("dashboard.title"),
		"active":     "dashboard",
		"user":       user,
		"teams":      teams,
		"jYear":      jYear,
		"jMonthName": helpers.JalaliMonthName(jMonth),
		"years":      helpers.JalaliYears(),
		"success":    ctx.Request().Query("success"),
		"error":      ctx.Request().Query("error"),
	})
}
