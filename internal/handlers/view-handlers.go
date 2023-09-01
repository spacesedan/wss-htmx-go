package handlers

import (
	"fmt"
	"log"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/kataras/blocks"
)

type ViewHandler struct {
	IndexView *blocks.Blocks
	LoginView *blocks.Blocks
}

func NewViewHandler() *ViewHandler {
	indexView := blocks.New("./views/index").Reload(true)
	_ = indexView.Load()

	loginView := blocks.New("./views/login").Reload(true)
	_ = loginView.Load()

	return &ViewHandler{
		IndexView: indexView,
		LoginView: loginView,
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
	fmt.Println(username)
	vars := map[string]interface{}{
		"Username": username.Value,
	}
	err := renderPage(w, v.IndexView, "index", "main", vars)
	if err != nil {
		log.Println(err)
		return
	}
}

func (v *ViewHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, v.LoginView, "index", "main", nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func renderPage(w http.ResponseWriter, block *blocks.Blocks, tmplName, layoutName string, data interface{}) error {
	err := block.ExecuteTemplate(w, tmplName, layoutName, data)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
