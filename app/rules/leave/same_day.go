package leave

import (
	"time"

	"github.com/goravel/framework/contracts/validation"
)

type SameDay struct {
}

// Signature The name of the rule.
func (receiver *SameDay) Signature() string {
	return "leave_same_day"
}

// Passes Determine if the validation rule passes.
func (receiver *SameDay) Passes(data validation.Data, val any, options ...any) bool {
	startRaw, _ := data.Get("start_at")
	endRaw, _ := data.Get("end_at")

	startStr, ok1 := startRaw.(string)
	endStr, ok2 := endRaw.(string)
	if !ok1 || !ok2 {
		return false
	}

	layout := "2006-01-02 15:04:05"
	start, err1 := time.Parse(layout, startStr)
	end, err2 := time.Parse(layout, endStr)
	if err1 != nil || err2 != nil {
		return false
	}

	// مقایسه فقط تاریخ (سال، ماه، روز)
	sy, sm, sd := start.Date()
	ey, em, ed := end.Date()

	return sy == ey && sm == em && sd == ed
}

// Message Get the validation error message.
func (receiver *SameDay) Message() string {
	return "امکان ثبت مرخصی در تاریخ های متفاوت امکان پذیر نیست."
}
