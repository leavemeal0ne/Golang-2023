package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

// RenderTemplate renders a template
func RenderTemplate(w http.ResponseWriter, tmpl string, data any) {
	// create a template cache
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = createTemplateCache()
	}

	var err error
	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)

	err = t.Execute(buf, data)
	if err != nil {
		log.Println(err)
	}

	// render the template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}
func parseDate(date time.Time) time.Time {
	result, _ := time.Parse("2006-01-02", date.String())
	return result
}

func createTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./ui/html/*_page.html")
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		//get only file name
		name := filepath.Base(page)
		ts, err := template.ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./ui/html/*_layout.html")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./ui/html/*_layout.html")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}

	return myCache, nil
}
