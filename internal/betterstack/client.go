package betterstack

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	betterStackBaseURL  = "https://uptime.betterstack.com"
	HeaderAuthorization = "Authorization"
	HeaderAccept        = "Accept"
	HeaderContentType   = "Content-Type"
	HeaderUserAgent     = "User-Agent"
	MediaTypeJSON       = "application/json"
)

type Client struct {
	httpClient *http.Client
	config     Config
}

type Config struct {
	APIToken string
	PageSize int
}

type Monitor struct {
	ID                string
	URL               string
	PronounceableName string
	Status            string
}

// Maps the relevant subset of fields from the 'list monitors' API
type MonitorListResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes *struct {
			URL               string `json:"url"`
			PronounceableName string `json:"pronounceable_name"`
			Status            string `json:"status"`
		} `json:"attributes"`
	} `json:"data"`
	Pagination *struct {
		First string `json:"first"`
		Last  string `json:"last"`
		Prev  string `json:"prev"`
		Next  string `json:"next"`
	} `json:"pagination"`
}

func NewClient(config Config) Client {
	if config.PageSize < 1 {
		config.PageSize = 50 // default https://betterstack.com/docs/uptime/api/pagination/
	}
	if config.PageSize > 250 {
		config.PageSize = 250 // maximum https://betterstack.com/docs/uptime/api/pagination/
	}
	return Client{
		config:     config,
		httpClient: &http.Client{Timeout: time.Duration(5) * time.Minute},
	}
}

func (c Client) execRequest(req *http.Request, expectedStatus int) (*http.Response, error) {
	req.Header.Set(HeaderAuthorization, "Bearer "+c.config.APIToken)
	req.Header.Set(HeaderAccept, MediaTypeJSON)
	req.Header.Set(HeaderContentType, MediaTypeJSON)
	req.Header.Add(HeaderUserAgent, "betterstack-exporter")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != expectedStatus {
		defer resp.Body.Close()
		result, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("got status %d, expected %d. Body: %b", resp.StatusCode, expectedStatus, result)
	}
	return resp, nil // caller should close resp.Body!
}

func (c Client) listMonitors() (*MonitorListResponse, error) {
	// Make HTTP request to the list monitors URL
	url := fmt.Sprintf("%s/api/v2/monitors?per_page=%d", betterStackBaseURL, c.config.PageSize)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.execRequest(req, http.StatusOK)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response
	var monitors MonitorListResponse
	err = json.NewDecoder(resp.Body).Decode(&monitors)
	if err != nil {
		return nil, err
	}
	return &monitors, nil
}

func (c Client) ListMonitors() ([]Monitor, error) {
	result := []Monitor{}
	monitors, err := c.listMonitors()
	if err != nil {
		return []Monitor{}, err
	}
	for {
		for _, monitor := range monitors.Data {
			result = append(result, Monitor{
				ID:                monitor.ID,
				PronounceableName: monitor.Attributes.PronounceableName,
				URL:               monitor.Attributes.URL,
				Status:            monitor.Attributes.Status,
			})
		}
		if !monitors.hasNext() {
			break // exit infinite loop
		}
		monitors, err = monitors.next(c)
		if err != nil {
			return []Monitor{}, err
		}
	}
	return result, nil
}

func (m MonitorListResponse) hasNext() bool {
	return m.Pagination != nil && m.Pagination.Next != ""
}

func (m MonitorListResponse) next(client Client) (*MonitorListResponse, error) {
	if !m.hasNext() {
		return nil, nil
	}

	// Make HTTP request to the next URL
	req, err := http.NewRequest(http.MethodGet, m.Pagination.Next, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.execRequest(req, http.StatusOK)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response
	var nextPage MonitorListResponse
	err = json.NewDecoder(resp.Body).Decode(&nextPage)
	if err != nil {
		return nil, err
	}
	return &nextPage, nil
}
