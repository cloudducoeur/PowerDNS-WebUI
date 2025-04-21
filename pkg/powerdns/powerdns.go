package powerdns

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

type PowerDNSClient struct {
	BaseURL  string
	APIKey   string
	ServerID string
	Client   *http.Client
}

// NewPowerDNSClient creates a new instance of PowerDNSClient.
//
// Parameters:
// - baseURL: The base URL of the PowerDNS API.
// - apiKey: The API key used for authentication with the PowerDNS API.
// - serverID: The ID of the PowerDNS server to interact with.
//
// Returns:
// - A pointer to a new PowerDNSClient instance configured with the provided parameters.
func NewPowerDNSClient(baseURL, apiKey, serverID string) *PowerDNSClient {
	return &PowerDNSClient{
		BaseURL:  baseURL,
		APIKey:   apiKey,
		ServerID: serverID,
		Client:   &http.Client{},
	}
}

// FetchZones retrieves all DNS zones from the PowerDNS API.
//
// Behavior:
// - Sends a GET request to the PowerDNS API to fetch the list of zones.
// - For each zone, fetches its DNS records using the FetchRecords method.
//
// Returns:
// - A slice of Zone structs containing the retrieved zones and their records.
// - An error if the API request fails, the response status is not OK, or the response cannot be decoded.
func (c *PowerDNSClient) FetchZones() ([]Zone, error) {
	url := fmt.Sprintf("%s/api/v1/servers/%s/zones", c.BaseURL, c.ServerID)

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

	var zones []Zone
	if err := json.NewDecoder(resp.Body).Decode(&zones); err != nil {
		return nil, fmt.Errorf("error decoding zones: %v", err)
	}

	// Récupérer les enregistrements pour chaque zone
	for i, zone := range zones {
		records, err := c.FetchRecords(zone.ID)
		if err != nil {
			return nil, err
		}
		zones[i].Records = records
	}

	return zones, nil
}

// FetchRecords retrieves DNS records for a given zone.
//
// Parameters:
// - zoneID: The ID of the zone for which to fetch DNS records.
//
// Behavior:
// - Sends a GET request to the PowerDNS API to fetch the records of the specified zone.
// - Decodes the response into a slice of DNSRecord structs.
//
// Returns:
// - A slice of DNSRecord structs containing the retrieved records.
// - An error if the API request fails, the response status is not OK, or the response cannot be decoded.
func (c *PowerDNSClient) FetchRecords(zoneID string) ([]DNSRecord, error) {
	url := fmt.Sprintf("%s/api/v1/servers/%s/zones/%s", c.BaseURL, c.ServerID, zoneID)

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", resp.Status)
	}

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

	if err := json.NewDecoder(resp.Body).Decode(&zoneResponse); err != nil {
		return nil, fmt.Errorf("error decoding records: %v", err)
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

// makeRequest is a utility method for making HTTP requests.
//
// Parameters:
// - method: The HTTP method to use (e.g., "GET", "POST").
// - url: The URL to send the request to.
// - body: The request body (optional, can be nil).
//
// Behavior:
// - Creates an HTTP request with the specified method, URL, and body.
// - Adds the API key and content type headers to the request.
// - Sends the request using the HTTP client.
//
// Returns:
// - The HTTP response from the server.
// - An error if the request creation or execution fails.
func (c *PowerDNSClient) makeRequest(method, url string, body interface{}) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("X-API-Key", c.APIKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	return resp, nil
}
