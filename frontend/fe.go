package frontend

import (
	"html/template"
)

// The frontend package is intended to do front end operations
// such as but not limited to loading html templates
func LoadHTML(filename string) (*template.Template, error) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		//http.Error(&w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	return tmpl, err
}