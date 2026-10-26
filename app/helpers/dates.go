package helpers

import (
	"fmt"
	"time"

	ptime "github.com/yaa110/go-persian-calendar"
)

// CurrentJalali returns the current Jalali year and month (1-12).
func CurrentJalali() (int, int) {
	pt := ptime.New(time.Now())
	return pt.Year(), int(pt.Month())
}

// JalaliShort formats a time as a short Persian date: 1403/04/09
func JalaliShort(t time.Time) string {
	return ptime.New(t).Format("yyyy/MM/dd")
}

// JalaliLong formats a time as a long Persian date: شنبه 9 تیر 1403
func JalaliLong(t time.Time) string {
	return ptime.New(t).Format("E d MMM yyyy")
}

// Clock formats the time portion as HH:mm
func Clock(t time.Time) string {
	return t.Format("15:04")
}

// JalaliMonthName returns the Persian month name for a 1-12 month number.
func JalaliMonthName(month int) string {
	return ptime.Month(month).String()
}

// JalaliYearRange returns the Gregorian [from, to) bounds of a Jalali year.
func JalaliYearRange(year int) (time.Time, time.Time) {
	from := ptime.Date(year, ptime.Farvardin, 1, 0, 0, 0, 0, time.UTC).Time()
	to := ptime.Date(year+1, ptime.Farvardin, 1, 0, 0, 0, 0, time.UTC).Time()
	return from, to
}

// JalaliMonthRange returns the Gregorian [from, to) bounds of a Jalali month.
func JalaliMonthRange(year, month int) (time.Time, time.Time) {
	from := ptime.Date(year, ptime.Month(month), 1, 0, 0, 0, 0, time.UTC).Time()
	ny, nm := year, month+1
	if nm > 12 {
		nm = 1
		ny = year + 1
	}
	to := ptime.Date(ny, ptime.Month(nm), 1, 0, 0, 0, 0, time.UTC).Time()
	return from, to
}

// JalaliYears returns the current and a few previous Jalali years (for selects).
func JalaliYears() []int {
	y, _ := CurrentJalali()
	return []int{y, y - 1, y - 2, y - 3, y - 4}
}

// FormatHours formats a duration (in seconds) as "X ساعت Y دقیقه".
func FormatHours(seconds float64) string {
	total := int(seconds)
	h := total / 3600
	m := (total % 3600) / 60
	return fmt.Sprintf("%d ساعت %d دقیقه", h, m)
}
