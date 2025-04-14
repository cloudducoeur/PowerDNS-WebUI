package powerdns

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// DNSRecord représente un enregistrement DNS.
type DNSRecord struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	TTL      int    `json:"ttl"`
	Disabled bool   `json:"disabled"`
}

// Zone représente une zone DNS.
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

// PowerDNSClient est un client pour interagir avec l'API PowerDNS.
type PowerDNSClient struct {
	BaseURL  string
	APIKey   string
	ServerID string
	Client   *http.Client
}

// NewPowerDNSClient crée une nouvelle instance de PowerDNSClient.
func NewPowerDNSClient(baseURL, apiKey, serverID string) *PowerDNSClient {
	return &PowerDNSClient{
		BaseURL:  baseURL,
		APIKey:   apiKey,
		ServerID: serverID,
		Client:   &http.Client{},
	}
}

// FetchZones récupère toutes les zones DNS.
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

// FetchRecords récupère les enregistrements DNS pour une zone donnée.
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

// makeRequest est une méthode utilitaire pour effectuer des requêtes HTTP.
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
