package utils

import (
	"net/http"
	"html/template"
)

var templates *template.Template

// LoadTemplates : To load the templates at given path
func LoadTemplates(pattern string) {
	templates = template.Must(template.ParseGlob(pattern))
}

// ExecuteTemplate : To execute the selected template
func ExecuteTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl, data)
}
