package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Sunpacker/go-booking-app/internal/config"
	"github.com/Sunpacker/go-booking-app/internal/models"
	"github.com/justinas/nosurf"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

var functions = template.FuncMap{}
var templatesFormat = "./templates/%s.tmpl"

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(templateData *models.TemplateData, r *http.Request) *models.TemplateData {
	templateData.Flash = app.Session.PopString(r.Context(), "flash")
	templateData.Error = app.Session.PopString(r.Context(), "error")
	templateData.Warning = app.Session.PopString(r.Context(), "warning")
	templateData.CSRFToken = nosurf.Token(r)
	return templateData
}

// RenderTemplate renders templates using html
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, templateData *models.TemplateData) error {
	var templateCache map[string]*template.Template

	if app.UseCache {
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	templateFromCache, ok := templateCache[tmpl]
	if !ok {
		return errors.New(fmt.Sprintf("cannot get template '%s' from cache", tmpl))
	}

	templateBuffer := new(bytes.Buffer)
	templateData = AddDefaultData(templateData, r)
	_ = templateFromCache.Execute(templateBuffer, templateData)

	_, err := templateBuffer.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}

func getTemplateFilepathPattern(pattern string) string {
	return fmt.Sprintf(templatesFormat, pattern)
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
