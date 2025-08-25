package providers

import (
	"github.com/goravel/framework/contracts/foundation"
	"leave/app/interfaces"
	"leave/app/repositories"
)

type AppServiceProvider struct {
}

func (receiver *AppServiceProvider) Register(app foundation.Application) {
	app.Singleton((*interfaces.UserInterface)(nil), func(app foundation.Application) (any, error) {
		return repositories.NewUserRepository(), nil
	})

	app.Singleton((*interfaces.LeaveInterface)(nil), func(app foundation.Application) (any, error) {
		return repositories.NewLeaveRepository(), nil
	})
}

func (receiver *AppServiceProvider) Boot(app foundation.Application) {

}
