package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/thanhphuocnguyen/go-bookings-app/pkg/config"
	"github.com/thanhphuocnguyen/go-bookings-app/pkg/handlers"
	"github.com/thanhphuocnguyen/go-bookings-app/pkg/render"
)

const portNumber = ":8083"

var appConfig config.AppConfig

func main() {
	// This is the entry point of the application
	appConfig.InProduction = false
	appConfig.UseCache = false

	appConfig.Session = scs.New()
	appConfig.Session.Lifetime = 24 * time.Hour
	appConfig.Session.Cookie.Persist = true
	appConfig.Session.Cookie.SameSite = http.SameSiteLaxMode
	appConfig.Session.Cookie.Secure = appConfig.InProduction

	// Initialize the template cache
	cache, err := render.InitTemplateCache(&appConfig)
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	appConfig.TemplateCache = cache
	render.InitializeRender(&appConfig)

	// Initialize a new repository
	repo := handlers.NewRepo(&appConfig)
	handlers.InitRepo(repo)

	fmt.Printf("Starting application on port %s\n", portNumber)

	server := &http.Server{
		Addr:    portNumber,
		Handler: routes(),
	}

	err = server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
