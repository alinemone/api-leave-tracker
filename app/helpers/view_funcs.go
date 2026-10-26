package helpers

import (
	"context"
	"html/template"

	"github.com/goravel/framework/facades"
)

// TemplateFuncs returns the custom functions available inside .tmpl views.
func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		// trans looks up a translation key from lang/<locale>.json (Laravel-style).
		"trans": func(key string) string {
			return facades.Lang(context.Background()).Get(key)
		},
		// jalali formats a time as a short Persian (Jalali) date: 1403/04/09
		"jalali": JalaliShort,
		// jalaliLong formats a time as a long Persian date: شنبه 9 تیر 1403
		"jalaliLong": JalaliLong,
		// clock formats the time portion as HH:mm
		"clock": Clock,
	}
}
