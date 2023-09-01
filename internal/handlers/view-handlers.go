package handlers

import (
	"log"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/kataras/blocks"
)

var views = blocks.New("./views").
	Reload(true).
	LayoutDir("layouts")

type ViewHandler struct {
	Views *blocks.Blocks
}

func NewViewHandler() *ViewHandler {
	_ = views.Load()
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
	vars := make(map[string]interface{})
	vars["username"] = username.String()
	err := v.renderPage(w, "index", "main", vars)
	if err != nil {
		log.Println(err)
		return
	}
}

func (v *ViewHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	err := v.renderPage(w, "index", "login", nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func (v *ViewHandler) renderPage(w http.ResponseWriter, tmplName, layoutName string, data interface{}) error {
	err := v.Views.ExecuteTemplate(w, tmplName, layoutName, data)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
