package handlers

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

// RenderTemplate renders an HTML template with the provided data.
//
// Parameters:
// - w: The HTTP response writer used to send the rendered template to the client.
// - tmpl: The name of the template file to render (located in the "templates/" directory).
// - data: A TemplateData struct containing the data to populate the template.
//
// Behavior:
// - Attempts to parse the specified template file from the "templates/" directory.
// - If parsing fails, responds with an HTTP 500 error and the error message.
// - Executes the template with the provided data and writes the output to the response.
// - If execution fails, responds with an HTTP 500 error and the error message.
func RenderTemplate(w http.ResponseWriter, tmpl string, data TemplateData) {
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

// RenderError renders an error message using the "index.html" template.
//
// Parameters:
// - w: The HTTP response writer used to send the rendered error message to the client.
// - message: A custom error message to display.
// - err: The error object containing additional details about the error.
//
// Behavior:
// - Constructs a TemplateData struct with the error message and details.
// - Calls RenderTemplate to render the "index.html" template with the error data.
func RenderError(w http.ResponseWriter, message string, err error) {
	data := TemplateData{
		Error: fmt.Sprintf("%s: %v", message, err),
	}
	RenderTemplate(w, "index.html", data)
}
