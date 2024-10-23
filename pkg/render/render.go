package render

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/thanhphuocnguyen/go-bookings-app/pkg/config"
	"github.com/thanhphuocnguyen/go-bookings-app/pkg/models"
)

// Function is a map of functions that can be used in the template
var function = template.FuncMap{}

var app *config.AppConfig

func InitializeRender(appConfig *config.AppConfig) {
	app = appConfig
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data *models.TemplateData) {
	var t *template.Template
	var ok bool
	data = AddDefaultData(data, r)
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
		log.Println("Error getting template files", err)
		return cache, err
	}
	for _, file := range tmplFiles {
		name := filepath.Base(file)

		// template.New(name).Funcs(function) is used to create a new template with the given name and function map
		ts, err := template.New(name).Funcs(function).ParseFiles(file)
		if err != nil {
			log.Println("Error parsing template", err)
			return cache, err
		}

		layoutTmpl, err := filepath.Glob("./templates/*.layout.tmpl")

		if err != nil {
			log.Println("Error getting layout files", err)
			return cache, err
		}

		if len(layoutTmpl) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				log.Println("Error parsing layout template", err)
				return cache, err
			}
		}

		cache[name] = ts
	}

	return cache, nil
}
