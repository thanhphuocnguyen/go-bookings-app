package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/config"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/models"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/render"
)

var appConfig config.AppConfig

const (
	pathToTemplates = "./../../templates"
	layoutSuffix    = ".layout.tmpl"
	pageSuffix      = ".page.tmpl"
)

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	// Initialize the template cache
	templateCache, err := InitTemplateCache()
	if err != nil {
		log.Fatalln("Error parsing templates: ", err)
	}

	appConfig.TemplateCache = templateCache
	// This is the entry point of the application
	appConfig.InProduction = false
	appConfig.UseCache = true

	mailChan := make(chan models.MailData)
	appConfig.MailChan = mailChan

	defer close(mailChan)

	listenForMail()

	appConfig.Session = scs.New()
	appConfig.Session.Lifetime = 24 * time.Hour
	appConfig.Session.Cookie.Persist = true
	appConfig.Session.Cookie.SameSite = http.SameSiteLaxMode
	appConfig.Session.Cookie.Secure = appConfig.InProduction

	appConfig.InfoLog = *log.New(log.Writer(), "INFO\t", log.Ldate|log.Ltime)
	appConfig.ErrorLog = *log.New(log.Writer(), "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	render.InitializeRenderer(&appConfig)

	// Initialize a new repository
	repo := InitializeTestingRepository(&appConfig)
	SetRepository(repo)

	os.Exit(m.Run())
}

func listenForMail() {
	go func() {
		for {
			<-appConfig.MailChan
		}
	}()
}

func NoSurf(next http.Handler) http.Handler {
	csrf := nosurf.New(next)
	csrf.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   appConfig.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrf
}

func SessionLoad(next http.Handler) http.Handler {
	return appConfig.Session.LoadAndSave(next)
}

func getRoutes() http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// comment out the following line to disable the CSRF protection middleware for testing
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/contact", Repo.Contact)
	// mux.Get("/generals-quarters", Repo.Generals)
	// mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.CreateReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.SearchAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/user/login", Repo.ShowLogin)
	mux.Post("/user/login", Repo.PostLogin)
	mux.Get("/user/logout", Repo.Logout)
	mux.Get("/user/register", Repo.ShowRegistration)
	mux.Post("/user/register", Repo.PostRegistration)

	mux.Get("/admin/dashboard", Repo.AdminDashboard)
	mux.Get("/admin/reservations/{id}", Repo.AdminShowReservation)
	mux.Post("/admin/reservations/{id}", Repo.AdminEditReservation)
	mux.Get("/admin/reservations/{id}/delete", Repo.AdminDeleteReservation)
	mux.Get("/admin/reservations/{id}/processed", Repo.AdminProcessedReservation)
	mux.Get("/admin/reservations-new", Repo.AdminNewReservations)
	mux.Get("/admin/reservations-all", Repo.AdminAllReservations)
	mux.Get("/admin/reservations-calendar", Repo.AdminReservationsCalendar)
	mux.Post("/admin/reservations-calendar", Repo.AdminPostReservationsCalendar)
	return mux
}

// Function is a map of functions that can be used in the template
var function = template.FuncMap{
	"humanDate":  render.HumanDate,
	"formatDate": render.FormatDate,
	"iterate":    render.Iterate,
	"add":        render.Add,
}

func InitTemplateCache() (map[string]*template.Template, error) {
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
