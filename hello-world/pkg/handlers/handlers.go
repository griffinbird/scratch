package handlers

import (
	"net/http"

	"github.com/griffinbird/booking/pkg/config"
	"github.com/griffinbird/booking/pkg/render"
	"github.com/griffinbird/go-course/pkg/config"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepp(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}

}

// NewHandlers set the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home.page.tmpl")
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.page.tmpl")
}
