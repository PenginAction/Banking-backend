package utils

import (
	"html/template"
	"net/http"
)


func RenderTemplate(w http.ResponseWriter, tmplPath string, data interface{}) error {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}
