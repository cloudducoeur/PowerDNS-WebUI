package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/cloudducoeur/PowerDNS-WebUI/pkg/powerdns"
)

type TemplateData struct {
	Zones []powerdns.Zone
	Query string
	Error string
}

func renderTemplate(w http.ResponseWriter, tmpl string, data TemplateData) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderError(w http.ResponseWriter, message string, err error) {
	data := TemplateData{
		Error: fmt.Sprintf("%s: %v", message, err),
	}
	renderTemplate(w, "index.html", data)
}
