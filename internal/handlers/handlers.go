package handlers

import (
	"net/http"
	"strings"

	"github.com/cloudducoeur/PowerDNS-WebUI/pkg/powerdns"
)

var powerDNSClient *powerdns.PowerDNSClient

// SetPowerDNSClient sets the PowerDNS client for handlers.
func SetPowerDNSClient(client *powerdns.PowerDNSClient) {
	powerDNSClient = client
}

// ListZonesHandler handles HTTP requests to list DNS zones.
//
// This handler retrieves DNS zones from the PowerDNS API and optionally filters
// the zones based on query parameters provided in the request.
//
// Query Parameters:
// - q: A search query string used to filter DNS records within zones.
// - type: The type of search to perform. Possible values are:
//   - "name": Filters records by their name.
//   - "content": Filters records by their content.
//   - "type": Filters records by their type.
//   - "all" (default): Filters records by name, content, or type.
//
// Behavior:
//   - If no query (`q`) is provided, all zones are returned.
//   - If a query is provided, only zones containing records that match the query
//     and search type are returned.
//
// Response:
// - On success, renders the "index.html" template with the filtered zones and query data.
// - On failure, renders an error message if zones cannot be fetched.
//
// Parameters:
// - w: The HTTP response writer used to send the response.
// - r: The HTTP request containing query parameters.
//
// Dependencies:
// - Fetches zones using the `powerDNSClient.FetchZones` method.
// - Filters records using the `filterRecords` helper function.
// - Renders templates using the `RenderTemplate` function.
func ListZonesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	searchType := r.URL.Query().Get("type")

	zones, err := powerDNSClient.FetchZones()
	if err != nil {
		RenderError(w, "Error fetching zones", err)
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

	RenderTemplate(w, "index.html", data)
}

// filterRecords filters DNS records based on a query string and search type.
//
// Parameters:
// - records: A slice of DNS records to filter.
// - query: The search query string used to match records.
// - searchType: The type of search to perform. Possible values are:
//   - "name": Filters records by their name.
//   - "content": Filters records by their content.
//   - "type": Filters records by their type.
//   - "all" (default): Filters records by name, content, or type.
//
// Behavior:
//   - Converts the query string to lowercase for case-insensitive matching.
//   - Iterates through the provided records and checks if each record matches
//     the query based on the specified search type.
//   - If a record matches, it is added to the filtered results.
//
// Returns:
// - A slice of DNS records that match the query and search type.
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
