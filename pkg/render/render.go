package render

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/thanhphuocnguyen/go-bookings-app/pkg/config"
	"github.com/thanhphuocnguyen/go-bookings-app/pkg/models"
)

// Function is a map of functions that can be used in the template
var function = template.FuncMap{}

var app *config.AppConfig

func InitializeRender(appConfig *config.AppConfig) {
	app = appConfig
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data *models.TemplateData) {
	var t *template.Template
	var ok bool

	if app.UseCache {
		t, ok = app.TemplateCache[tmpl]
		if !ok {
			panic("Could not get template from cache")
		}
	} else {
		t, ok = app.TemplateCache[tmpl]
		if !ok {
			panic("Could not get template from cache")
		}
	}

	err := t.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func InitTemplateCache(appConfig *config.AppConfig) (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)
	tmplFiles, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return cache, err
	}
	for _, file := range tmplFiles {
		name := filepath.Base(file)

		// template.New(name).Funcs(function) is used to create a new template with the given name and function map
		ts, err := template.New(name).Funcs(function).ParseFiles(file)
		if err != nil {
			return cache, err
		}

		layoutTmpl, err := filepath.Glob("./templates/*.layout.tmpl")

		if err != nil {
			return cache, err
		}

		if len(layoutTmpl) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return cache, err
			}
		}

		cache[name] = ts
	}

	return cache, nil
}
