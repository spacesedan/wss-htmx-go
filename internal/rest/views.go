package rest

import (
	"html/template"
	"net/http"

	chi "github.com/go-chi/chi/v5"
)

type ViewHandler struct{}

func NewViewHandler() *ViewHandler {
	return &ViewHandler{}
}

func (v *ViewHandler) Register(r *chi.Mux) {
	r.Get("/", v.Index)
}

func (v *ViewHandler) Index(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles("templates/index.html"))
	_ = tmp.Execute(w, nil)
}
