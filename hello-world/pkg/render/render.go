package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

// NewTemplate sets the config for the template package
func NewTemplate(a *config.AppConfig) {
	app = a
}

// RenderTemplate renders a templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string) {
	var tc map[string]*template.Template
	
	
	if app.UseCache {}
		// get the template cache from the app configi
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal(" could not get a template from the template cache")
	}

	// create a template cache
	tc, err := createTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer) // buffer to put template

	err = t.Execute(buf, nil)
	if err != nil {
		log.Println(err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func createTemplateCache() (map[string]*template.Template, error) {
	// myCache := make(map[string]*.template.createTemplateCache)
	myCache := map[string]*template.Template{}

	// get all of the files named * .page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.pages.tmpl")
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl, ts = template set
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./template/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err

			}
		}

		myCache[name] = ts // name of file
	}

	return myCache, nil

}
