package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	UseCache      bool
	InfoLog       log.Logger
	ErrorLog      log.Logger
	TemplateCache map[string]*template.Template
	InProduction  bool
	Session       *scs.SessionManager
}
