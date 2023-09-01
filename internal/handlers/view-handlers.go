package handlers

import (
	"fmt"
	"log"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/kataras/blocks"
)

type ViewHandler struct {
	Views *blocks.Blocks
}

const HTML_HEADER = "text/html; charset=utf-8"

func NewViewHandler() *ViewHandler {
	views := blocks.New("./views").Reload(true)
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
	w.Header().Set("Content-Type", HTML_HEADER)

	// get the username from cookie value
	username, _ := r.Cookie("username")
	fmt.Println(username)

	vars := map[string]interface{}{
		"Username": username.Value,
	}
	err := v.renderPage(w, "index", "main", vars)
	if err != nil {
		log.Println(err)
		return
	}
}

func (v *ViewHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", HTML_HEADER)

	err := v.renderPage(w, "login", "main", nil)
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
