package leave

import (
	"github.com/goravel/framework/contracts/event"
)

type LeaveCreatedEvent struct {
}

func (receiver *LeaveCreatedEvent) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}
