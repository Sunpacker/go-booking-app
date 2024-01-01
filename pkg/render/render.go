package render

import (
	"bytes"
	"fmt"
	"github.com/Sunpacker/go-booking-app/pkg/config"
	"github.com/Sunpacker/go-booking-app/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var functions = template.FuncMap{}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(templateData *models.TemplateData) *models.TemplateData {
	return templateData
}

// RenderTemplate renders templates using html
func RenderTemplate(w http.ResponseWriter, tmpl string, templateData *models.TemplateData) {
	var templateCache map[string]*template.Template

	if app.UseCache {
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	templateFromCache, ok := templateCache[tmpl]
	if !ok {
		log.Fatal("template not found")
	}

	templateBuffer := new(bytes.Buffer)
	templateData = AddDefaultData(templateData)
	_ = templateFromCache.Execute(templateBuffer, templateData)

	_, err := templateBuffer.WriteTo(w)
	if err != nil {
		fmt.Println("error while writing template to browser:", err)
		return
	}
}

func getTemplateFilepathPattern(pattern string) string {
	return "./templates/" + pattern + ".tmpl"
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	templateCache := map[string]*template.Template{}

	pagePattern := getTemplateFilepathPattern("*.page")
	pages, err := filepath.Glob(pagePattern)
	if err != nil {
		return templateCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		templateSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		layoutPattern := getTemplateFilepathPattern("*.layout")
		matches, err := filepath.Glob(layoutPattern)
		if err != nil {
			return templateCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob(layoutPattern)
			if err != nil {
				return templateCache, err
			}
		}

		pageName := strings.Split(name, ".")[0]
		templateCache[pageName] = templateSet
	}

	return templateCache, nil
}
