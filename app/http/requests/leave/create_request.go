package leave

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
	ptime "github.com/yaa110/go-persian-calendar"

	"leave/app/helpers"
)

type CreateRequest struct {
	Type    string    `form:"type" json:"type" example:"daily"`
	StartAt time.Time `form:"start_at" json:"start_at" example:"2025-08-11 14:00:00"`
	EndAt   time.Time `form:"end_at" json:"end_at" example:"2025-08-11 14:00:00"`
}

func (r *CreateRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *CreateRequest) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CreateRequest) Rules(ctx http.Context) map[string]string {
	user, _ := helpers.CurrentUser(ctx)

	return map[string]string{
		"type":     "required",
		"start_at": "required",
		"end_at":   fmt.Sprintf("required|leave_overlap:%d|leave_same_day", user.ID),
	}
}

func (r *CreateRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CreateRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

// PrepareForValidation converts incoming Jalali dates ("1403/04/09" or
// "1403/04/09 09:00:00") to Gregorian "2006-01-02 15:04:05" before binding.
func (r *CreateRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	for _, key := range []string{"start_at", "end_at"} {
		raw, ok := data.Get(key)
		if !ok {
			continue
		}
		val, ok := raw.(string)
		if !ok || val == "" {
			continue
		}
		if converted, ok := jalaliToGregorian(val); ok {
			data.Set(key, converted)
		}
	}
	return nil
}

// jalaliToGregorian parses a Jalali date(+time) string and returns the
// Gregorian equivalent formatted as "2006-01-02 15:04:05".
func jalaliToGregorian(value string) (string, bool) {
	value = toEnglishDigits(strings.TrimSpace(value))
	if value == "" || !strings.Contains(value, "/") {
		return "", false
	}

	datePart := value
	timePart := ""
	if idx := strings.Index(value, " "); idx != -1 {
		datePart = value[:idx]
		timePart = strings.TrimSpace(value[idx+1:])
	}

	dp := strings.Split(datePart, "/")
	if len(dp) != 3 {
		return "", false
	}
	year, e1 := strconv.Atoi(strings.TrimSpace(dp[0]))
	month, e2 := strconv.Atoi(strings.TrimSpace(dp[1]))
	day, e3 := strconv.Atoi(strings.TrimSpace(dp[2]))
	if e1 != nil || e2 != nil || e3 != nil {
		return "", false
	}

	hour, minute, sec := 0, 0, 0
	if timePart != "" {
		tp := strings.Split(timePart, ":")
		if len(tp) >= 1 {
			hour, _ = strconv.Atoi(tp[0])
		}
		if len(tp) >= 2 {
			minute, _ = strconv.Atoi(tp[1])
		}
		if len(tp) >= 3 {
			sec, _ = strconv.Atoi(tp[2])
		}
	}

	g := ptime.Date(year, ptime.Month(month), day, hour, minute, sec, 0, time.UTC).Time()
	return g.Format("2006-01-02 15:04:05"), true
}

// toEnglishDigits normalizes Persian/Arabic digits to ASCII.
func toEnglishDigits(s string) string {
	replacer := strings.NewReplacer(
		"۰", "0", "۱", "1", "۲", "2", "۳", "3", "۴", "4",
		"۵", "5", "۶", "6", "۷", "7", "۸", "8", "۹", "9",
		"٠", "0", "١", "1", "٢", "2", "٣", "3", "٤", "4",
		"٥", "5", "٦", "6", "٧", "7", "٨", "8", "٩", "9",
	)
	return replacer.Replace(s)
}
