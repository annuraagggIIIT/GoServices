// pkg/render/render.go
package render

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"bytes"
	"github.com/annuraagggIIIT/Go-Practice/config"
	"github.com/annuraagggIIIT/Go-Practice/models"
)

var app *config.AppConfig

// NewTemplates sets the config for template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	//get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td)
	_ = t.Execute(buf, td)

	// render template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("error writing template to browser:", err)
	}
}
// pkg/render/render.go
// ...

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}





	// get all files named *page.tmpl from the ./templates folder
	pages, err := filepath.Glob("./templates/*page.tmpl")
	if err != nil {
		log.Println("Error reading page templates:", err)
		return myCache, err
	}
	
	log.Printf("Found %d page templates: %v", len(pages), pages)

	// range through all files ending with *page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		log.Printf("Processing template: %s", name)

		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			log.Println("Error parsing template:", err)
			return myCache, err
		}

		layouts, err := filepath.Glob("./templates/*layout.tmpl")
		if err != nil {
			log.Println("Error reading layout templates:", err)
			return myCache, err
		}

		if len(layouts) > 0 {
			ts, err = ts.ParseGlob("./templates/*layout.tmpl")
			if err != nil {
				log.Println("Error parsing layout templates:", err)
				return myCache, err
			}
		}

		myCache[name] = ts
		log.Println("Successfully parsed template:", name)
	}

	return myCache, nil
}


// var tc = make(map[string]*template.Template)

// func RenderTemplate(w http.ResponseWriter, t string) {
// 	var tmpl *template.Template
// 	var err error

// 	_, inMap := tc[t]
// 	if !inMap {
// 		// need to create template
// 		log.Println("creating template and adding to cache")
// 		err = createTemplateCache(t)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	} else {
// 		log.Println("using cached template")
// 	}

// 	tmpl = tc[t]
// 	err = tmpl.Execute(w, nil)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// func createTemplateCache(t string) error {
// 	templates := []string{
// 		fmt.Sprintf("./templates/%s", t),
// 		"./templates/base.layout.tmpl",
// 	}

// 	tmpl, err := template.ParseFiles(templates...)
// 	if err != nil {
// 		return err
// 	}

// 	tc[t] = tmpl
// 	return nil
// }