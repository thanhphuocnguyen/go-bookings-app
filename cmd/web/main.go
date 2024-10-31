package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/config"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/driver"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/handlers"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/helpers"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/models"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/render"
)

const portNumber = ":8083"

var appConfig config.AppConfig

func main() {
	db, err := run()

	if err != nil {
		log.Fatalln("Error starting application: ", err)
	}
	defer db.SQL.Close()
	defer close(appConfig.MailChan)

	listenForMail()

	if err != nil {
		log.Fatalln("Error starting application: ", err)
	}

	fmt.Printf("Starting application on port %s\n", portNumber)
	server := &http.Server{
		Addr:    portNumber,
		Handler: routes(),
	}
	err = server.ListenAndServe()
	log.Fatalln("Error starting server: ", err)
}

func run() (*driver.DB, error) {
	gob.Register(models.Reservation{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	// Initialize the template cache
	templateCache, err := render.InitializeTmplCache()
	if err != nil {
		log.Fatalln("Error parsing templates: ", err)
		return nil, err
	}

	mailChan := make(chan models.MailData)
	appConfig.MailChan = mailChan

	appConfig.TemplateCache = templateCache
	// This is the entry point of the application
	appConfig.InProduction = false
	appConfig.UseCache = false

	appConfig.InfoLog = *log.New(log.Writer(), "INFO\t", log.Ldate|log.Ltime)
	appConfig.ErrorLog = *log.New(log.Writer(), "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	appConfig.Session = scs.New()
	appConfig.Session.Lifetime = 24 * time.Hour
	appConfig.Session.Cookie.Persist = true
	appConfig.Session.Cookie.SameSite = http.SameSiteLaxMode
	appConfig.Session.Cookie.Secure = appConfig.InProduction

	// Initialize Database
	db, err := driver.InitializeDatabase("host=localhost port=5432 dbname=go_bookings user=postgres password=postgres sslmode=disable")

	if err != nil {
		log.Fatalln("Cannot connect to database: ", err)
		return nil, err
	}

	helpers.InitHelper(&appConfig)
	render.InitializeRenderer(&appConfig)
	// Initialize a new repository
	repo := handlers.InitializeRepository(&appConfig, db)
	handlers.SetRepository(repo)
	return db, nil
}
