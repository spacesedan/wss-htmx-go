package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type RestHandler struct{}

func NewRestHandler() *RestHandler {
	return &RestHandler{}
}

func (h *RestHandler) Register(r *chi.Mux) {
	r.Post("/username", h.Username)
}

func (h *RestHandler) Username(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	userName := r.Form.Get("username")
	cookie := http.Cookie{
		Name:     "username",
		Value:    userName,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &cookie)

	w.Header().Set("HX-Redirect", "http://localhost:8080/")
}
