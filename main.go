package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	PowerDNSURL string `json:"powerdns_url"`
	APIKey      string `json:"api_key"`
	ServerID    string `json:"server_id"`
}

type DNSRecord struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	TTL      int    `json:"ttl"`
	Disabled bool   `json:"disabled"`
}

type Zone struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	Records        []DNSRecord `json:"records"`
	Serial         int         `json:"serial"`
	Kind           string      `json:"kind"`
	DNSSec         bool        `json:"dnssec"`
	Account        string      `json:"account"`
	LastCheck      int         `json:"last_check"`
	NotifiedSerial int         `json:"notified_serial"`
}

type TemplateData struct {
	Zones []Zone
	Query string
	Error string
}

var config Config

func main() {
	loadConfig()

	http.HandleFunc("/", listZonesHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Starting the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func loadConfig() {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error reading configuration file: %v", err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Error parsing configuration file: %v", err)
	}
}

func listZonesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	searchType := r.URL.Query().Get("type")

	zones, err := fetchZones()
	if err != nil {
		renderError(w, "Error fetching zones", err)
		return
	}

	var filteredZones []Zone
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

func filterRecords(records []DNSRecord, query, searchType string) []DNSRecord {
	var filtered []DNSRecord
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

func fetchZones() ([]Zone, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/v1/servers/%s/zones", config.PowerDNSURL, config.ServerID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", config.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	var zones []Zone
	err = json.Unmarshal(body, &zones)
	if err != nil {
		return nil, err
	}

	// Fetch records for each zone
	for i, zone := range zones {
		records, err := fetchRecords(zone.ID)
		if err != nil {
			return nil, err
		}
		zones[i].Records = records
	}

	return zones, nil
}

func fetchRecords(zoneID string) ([]DNSRecord, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/api/v1/servers/%s/zones/%s", config.PowerDNSURL, config.ServerID, zoneID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("X-API-Key", config.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	// Structure corresponding to the PowerDNS response
	var zoneResponse struct {
		Name   string `json:"name"`
		RRsets []struct {
			Name    string `json:"name"`
			Type    string `json:"type"`
			TTL     int    `json:"ttl"`
			Records []struct {
				Content  string `json:"content"`
				Disabled bool   `json:"disabled"`
			} `json:"records"`
		} `json:"rrsets"`
	}

	if err := json.Unmarshal(body, &zoneResponse); err != nil {
		return nil, fmt.Errorf("JSON decoding error: %v - Response: %s", err, string(body))
	}

	var records []DNSRecord
	for _, rrset := range zoneResponse.RRsets {
		for _, record := range rrset.Records {
			records = append(records, DNSRecord{
				Name:     rrset.Name,
				Type:     rrset.Type,
				Content:  record.Content,
				TTL:      rrset.TTL,
				Disabled: record.Disabled,
			})
		}
	}

	return records, nil
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
