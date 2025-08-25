package leave

import (
	"errors"
	"fmt"
	"leave/app/models"
	"leave/services/gap"

	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"
	ptime "github.com/yaa110/go-persian-calendar"
)

type LeaveDeletedNotificationListener struct {
}

func (receiver *LeaveDeletedNotificationListener) Signature() string {
	return "leave_deleted_notification_listener"
}

func (receiver *LeaveDeletedNotificationListener) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     false,
		Connection: "",
		Queue:      "",
	}
}

func (receiver *LeaveDeletedNotificationListener) Handle(args ...any) error {
	//gepEnabled := facades.Config().GetBool("GAP_ENABLED")
	//if gepEnabled == false {
	//	return nil
	//}

	UserName, _ := args[0].(string)
	ID, _ := args[1].(uint)

	var leave models.Leave

	if err := facades.Orm().Query().WithTrashed().Find(&leave, "id=?", ID); err != nil {
		return err
	}

	if leave.ID == 0 {
		return errors.New("leave not found")
	}

	var message string
	StartAt := ptime.New(leave.StartAt)
	message = fmt.Sprintf(
		"مرخصی %s در تاریخ %s %d %s حذف شد.",
		UserName,
		StartAt.Weekday(),
		StartAt.Day(),
		StartAt.Month().String(),
	)

	gapClient := gap.NewClient()
	err := gapClient.SendMessage(message, "")

	if err != nil {
		panic(err)
	}

	return nil

}
