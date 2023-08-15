package handlers

import (
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	chi "github.com/go-chi/chi/v5"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./templates"),
	jet.InDevelopmentMode(),
)

type ViewHandler struct {
	Views *jet.Set
}

func NewViewHandler() *ViewHandler {
	return &ViewHandler{
		Views: views,
	}
}

func (v *ViewHandler) Register(r *chi.Mux) {
	r.Route("/", func(r chi.Router) {
		r.Use(CheckUsernameCookie)
		r.Get("/", v.Index)
	})
	r.Get("/login", v.LoginPage)
}

func (v *ViewHandler) Index(w http.ResponseWriter, r *http.Request) {
	username, _ := r.Cookie("username")
	vars := make(jet.VarMap)
	vars.Set("user", username.Value)
	err := v.renderPage(w, "index.jet", vars)
	if err != nil {
		log.Println(err)
		return
	}
}

func (v *ViewHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	err := v.renderPage(w, "login.jet", nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func (v *ViewHandler) renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := v.Views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
