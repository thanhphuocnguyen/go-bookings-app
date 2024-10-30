package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/config"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/models"
)

// Function is a map of functions that can be used in the template
var function = template.FuncMap{}

var app *config.AppConfig

func InitializeRenderer(appConfig *config.AppConfig) {
	app = appConfig
}

var (
	pathToTemplates = "./templates"
	layoutSuffix    = ".layout.tmpl"
	pageSuffix      = ".page.tmpl"
)

func ApplyDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	return td
}

func RenderTmpl(w http.ResponseWriter, r *http.Request, tmpl string, data *models.TemplateData) error {
	var templateCache map[string]*template.Template
	if app.UseCache {
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = InitializeTmplCache()
	}
	buffer := new(bytes.Buffer)
	if tmpl, ok := templateCache[tmpl]; ok {
		data = ApplyDefaultData(data, r)
		err := tmpl.Execute(buffer, data)
		if err != nil {
			log.Println("Error executing template: ", err)
			return err
		}
		_, err = buffer.WriteTo(w)
		if err != nil {
			log.Println("Error writing template to browser: ", err)
			return err
		}
	} else {
		return fmt.Errorf("template not found in cache")
	}
	return nil
}

func InitializeTmplCache() (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)
	tmplFiles, err := filepath.Glob(fmt.Sprintf("%s/*%s", pathToTemplates, pageSuffix))
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

		layoutTmpl, err := filepath.Glob(fmt.Sprintf("%s/*%s", pathToTemplates, layoutSuffix))

		if err != nil {
			log.Println("Error getting layout files", err)
			return cache, err
		}

		if len(layoutTmpl) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*%s", pathToTemplates, layoutSuffix))
			if err != nil {
				log.Println("Error parsing layout template", err)
				return cache, err
			}
		}

		cache[name] = ts
	}
	return cache, nil
}
