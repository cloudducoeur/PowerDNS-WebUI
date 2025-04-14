package main

import (
	"net/http"
	"strings"

	"github.com/cloudducoeur/PowerDNS-WebUI/pkg/powerdns"
)

func listZonesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	searchType := r.URL.Query().Get("type")

	zones, err := powerDNSClient.FetchZones()
	if err != nil {
		renderError(w, "Error fetching zones", err)
		return
	}

	var filteredZones []powerdns.Zone
	if query != "" {
		for _, zone := range zones {
			filteredRecords := filterRecords(zone.Records, query, searchType)
			if len(filteredRecords) > 0 {
				filteredZone := zone
				filteredZone.Records = filteredRecords
				filteredZones = append(filteredZones, filteredZone)
			}
		}
	} else {
		filteredZones = zones
	}

	data := TemplateData{
		Zones: filteredZones,
		Query: query,
	}

	renderTemplate(w, "index.html", data)
}

func filterRecords(records []powerdns.DNSRecord, query, searchType string) []powerdns.DNSRecord {
	var filtered []powerdns.DNSRecord
	query = strings.ToLower(query)

	for _, record := range records {
		var match bool

		switch searchType {
		case "name":
			match = strings.Contains(strings.ToLower(record.Name), query)
		case "content":
			match = strings.Contains(strings.ToLower(record.Content), query)
		case "type":
			match = strings.Contains(strings.ToLower(record.Type), query)
		default: // "all"
			match = strings.Contains(strings.ToLower(record.Name), query) ||
				strings.Contains(strings.ToLower(record.Content), query) ||
				strings.Contains(strings.ToLower(record.Type), query)
		}

		if match {
			filtered = append(filtered, record)
		}
	}

	return filtered
}
