package render

import (
	"net/http"
	"testing"

	"github.com/thanhphuocnguyen/go-bookings-app/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var data models.TemplateData
	request, err := getSession()
	if err != nil {
		t.Error(err)
	}
	app.Session.Put(request.Context(), "flash", "123")
	result := ApplyDefaultData(&data, request)

	if result == nil {
		t.Error("Expected a pointer to TemplateData, but got nil")
	} else if result.Flash != "123" {
		t.Error("Expected 123, but got", result.Flash)
	}
}

type myWriter struct{}

func (mw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (mw *myWriter) WriteHeader(i int) {}

func (mw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}

func TestRenderTemplate(t *testing.T) {
	tc, err := InitializeTmplCache()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	err = RenderTmpl(&ww, r, "home.page.tmpl", &models.TemplateData{})

	if err != nil {
		t.Error("Error writing template to browser: ", err)
	}

	err = RenderTmpl(&ww, r, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("Rendered template that does not exist")
	}

}

func getSession() (*http.Request, error) {
	request, err := http.NewRequest("GET", "/testing", nil)
	if err != nil {
		return nil, err
	}

	ctx := request.Context()
	ctx, _ = app.Session.Load(ctx, request.Header.Get("X-Session"))
	request = request.WithContext(ctx)
	return request, nil
}
