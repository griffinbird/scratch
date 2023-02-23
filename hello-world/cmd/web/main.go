package main

import (
	"fmt"

	"github.com/griffinbird/bookings/pkg/config"
	"github.com/griffinbird/bookings/pkg/handlers"
	"github.com/griffinbird/bookings/pkg/render"
	"github.com/griffinbird/go-course/pkg/handlers"

	"log"
	"net/http"
)

const portNumber = ":8080"

// main is the main function
func main() {
	var app config.AppConfig

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepp(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	http.HandleFunc("/", handlers.Repo.Home)
	http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))
	_ = http.ListenAndServe(portNumber, nil)
}
