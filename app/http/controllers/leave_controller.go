package controllers

import (
	leaveEvent "leave/app/events/leave"
	"leave/app/helpers"
	leaveRequest "leave/app/http/requests/leave"
	"leave/app/http/responses"
	response "leave/app/http/responses/leave"
	"leave/app/interfaces"
	"leave/app/repositories"

	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type LeaveController struct {
	repository interfaces.LeaveInterface
	// Dependent services
}

func NewLeaveController() *LeaveController {
	return &LeaveController{
		repository: repositories.NewLeaveRepository(),
	}
}

func (r *LeaveController) Report(ctx http.Context) http.Response {
	page := ctx.Request().QueryInt("page", 1)
	perPage := ctx.Request().QueryInt("per_page", 10)

	leaves, total, _ := r.repository.Report(page, perPage)

	return ctx.Response().Success().Json(
		responses.NewPagination(leaves, page, perPage, total),
	)
}

func (r *LeaveController) Store(ctx http.Context) http.Response {
	var req leaveRequest.CreateRequest

	errors, _ := ctx.Request().ValidateRequest(&req)

	if errors != nil {
		return ctx.Response().Json(
			http.StatusUnprocessableEntity,
			responses.NewErrorResponse(errors),
		)
	}

	user, err := helpers.CurrentUser(ctx)

	if err != nil {
		return ctx.Response().Json(
			http.StatusUnauthorized,
			responses.NewErrorResponse("Unauthorized"),
		)
	}

	leave, err := r.repository.Create(user.ID, req.Type, req.StartAt, req.EndAt)

	if err != nil || leave.ID == 0 {
		return ctx.Response().Json(
			http.StatusInternalServerError,
			responses.NewErrorResponse("create leave error"),
		)
	}

	_ = facades.Event().Job(&leaveEvent.LeaveCreatedEvent{}, []event.Arg{
		{Type: "string", Value: user.Name},
		{Type: "string", Value: leave.StartAt.String()},
		{Type: "string", Value: leave.EndAt.String()},
		{Type: "string", Value: leave.Type},
	}).Dispatch()

	return ctx.Response().Success().Json(leave)
}

func (r *LeaveController) ListUserLeaves(ctx http.Context) http.Response {
	user, _ := helpers.CurrentUser(ctx)
	page := ctx.Request().QueryInt("page", 1)
	perPage := ctx.Request().QueryInt("per_page", 10)

	leaves, total, _ := r.repository.GetLeavesByUserID(user.ID, page, perPage)

	result := responses.NewPagination(
		response.NewLeaveCollection(leaves),
		page,
		perPage,
		total,
	)

	return ctx.Response().Success().Json(result)
}

func (r *LeaveController) Delete(ctx http.Context) http.Response {
	var req leaveRequest.DeleteRequest

	errors, _ := ctx.Request().ValidateRequest(&req)

	if errors != nil {
		return ctx.Response().Json(
			http.StatusUnprocessableEntity,
			responses.NewErrorResponse(errors),
		)
	}

	currentUser, _ := helpers.CurrentUser(ctx)

	err := r.repository.Delete(currentUser.ID, req.ID)

	if err != nil {
		return ctx.Response().Json(
			http.StatusInternalServerError,
			responses.NewErrorResponse(err),
		)
	}

	_ = facades.Event().Job(&leaveEvent.LeaveDeletedEvent{}, []event.Arg{
		{Type: "string", Value: currentUser.Name},
		{Type: "string", Value: req.ID},
	}).Dispatch()

	return ctx.Response().Success().Json(
		responses.NewPublicResponse(true),
	)
}
