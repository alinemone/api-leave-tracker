package leave

import (
	"fmt"
	"leave/services/gap"
	"time"

	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/facades"
	ptime "github.com/yaa110/go-persian-calendar"
)

type LeaveCreatedNotificationListener struct {
}

func (receiver *LeaveCreatedNotificationListener) Signature() string {
	return "send_notification_gap_listener"
}

func (receiver *LeaveCreatedNotificationListener) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     false,
		Connection: "",
		Queue:      "",
	}
}

func (receiver *LeaveCreatedNotificationListener) Handle(args ...any) error {

	gepEnabled := facades.Config().GetBool("GAP_ENABLED")
	if gepEnabled == false {
		return nil
	}

	UserName, _ := args[0].(string)
	StartAtStr, _ := args[1].(string)
	EndAtStr, _ := args[2].(string)
	Type, _ := args[3].(string)

	layout := "2006-01-02 15:04:05 -0700 -0700"
	startAtTime, _ := time.Parse(layout, StartAtStr)
	endAtTime, _ := time.Parse(layout, EndAtStr)

	StartAt := ptime.New(startAtTime)
	EndAt := ptime.New(endAtTime)

	var message string

	if Type == "daily" {
		dayName := StartAt.Weekday()
		message = fmt.Sprintf(
			"%s در تاریخ %s %d %s مرخصی ثبت کرد.",
			UserName,
			dayName,
			StartAt.Day(),
			StartAt.Month().String(),
		)

	} else if Type == "hourly" {
		dayName := StartAt.Weekday()
		message = fmt.Sprintf(
			"%s در تاریخ %s %d %s از ساعت %02d:%02d تا ساعت %02d:%02d مرخصی ساعتی ثبت کرد. ⏰",
			UserName,
			dayName,
			StartAt.Day(),
			StartAt.Month().String(),
			StartAt.Hour(),
			StartAt.Minute(),
			EndAt.Hour(),
			EndAt.Minute(),
		)

	} else {
		message = fmt.Sprintf("%s مرخصی ثبت کرد", UserName)
	}

	gapClient := gap.NewClient()
	err := gapClient.SendMessage(message, "")

	if err != nil {
		panic(err)
	}

	return nil
}
