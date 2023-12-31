package render

import (
	"fmt"
	"html/template"
	"net/http"
)

// RenderTemplate renders templates using html
func RenderTemplate(w http.ResponseWriter, tmpl string) {
	parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl + ".tmpl")
	err := parsedTemplate.Execute(w, nil)

	if err != nil {
		fmt.Println("error while parsing template:", err)
		return
	}
}
