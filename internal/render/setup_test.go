package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/config"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/models"
)

var testApp config.AppConfig

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})
	pathToTemplates = "./../../templates"

	// This is the entry point of the application
	testApp.InProduction = false
	testApp.UseCache = true

	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction

	testApp.Session = session
	app = &testApp
	os.Exit(m.Run())
}
