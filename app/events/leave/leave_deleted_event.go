package leave

import "github.com/goravel/framework/contracts/event"

type LeaveDeletedEvent struct {
}

func (receiver *LeaveDeletedEvent) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}
