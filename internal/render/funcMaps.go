package render

import (
	"html/template"
	"time"
)

// Function is a map of functions that can be used in the template
var function = template.FuncMap{
	"humanDate":  HumanDate,
	"formatDate": FormatDate,
	"iterate":    Iterate,
	"add":        Add,
}

func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func FormatDate(t time.Time, f string) string {
	return t.Format(f)
}

func Iterate(cnt int) []int {
	var i int
	var items []int
	for i = 0; i < cnt; i++ {
		items = append(items, i)
	}
	return items
}

func Add(a, b int) int {
	return a + b
}
