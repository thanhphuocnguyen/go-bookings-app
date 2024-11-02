package main

import (
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/thanhphuocnguyen/go-bookings-app/internal/helpers"
)

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

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(w, r) {
			appConfig.Session.Put(r.Context(), "error", "Log in first")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
