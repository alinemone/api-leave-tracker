package helpers

import "strings"

var digitReplacer = strings.NewReplacer(
	"۰", "0", "۱", "1", "۲", "2", "۳", "3", "۴", "4",
	"۵", "5", "۶", "6", "۷", "7", "۸", "8", "۹", "9",
	"٠", "0", "١", "1", "٢", "2", "٣", "3", "٤", "4",
	"٥", "5", "٦", "6", "٧", "7", "٨", "8", "٩", "9",
)

// EnglishDigits converts Persian/Arabic digits in s to ASCII 0-9 (used to
// normalize mobile numbers so they are always stored/compared in English).
func EnglishDigits(s string) string {
	return digitReplacer.Replace(s)
}
