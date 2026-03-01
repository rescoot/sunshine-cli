package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rescoot/sunshine-cli/internal/auth"
	"github.com/rescoot/sunshine-cli/internal/config"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	cfg        *config.Config
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		baseURL:    strings.TrimRight(cfg.Server, "/") + "/api/v1",
		httpClient: &http.Client{},
		cfg:        cfg,
	}
}

func (c *Client) get(path string, result interface{}) error {
	return c.do("GET", path, nil, result)
}

func (c *Client) post(path string, body interface{}, result interface{}) error {
	return c.do("POST", path, body, result)
}

func (c *Client) put(path string, body interface{}, result interface{}) error {
	return c.do("PUT", path, body, result)
}

func (c *Client) delete(path string, result interface{}) error {
	return c.do("DELETE", path, nil, result)
}

func (c *Client) do(method, path string, body interface{}, result interface{}) error {
	tokens, err := c.getValidTokens()
	if err != nil {
		return fmt.Errorf("authentication required: run 'sunshine auth login' first")
	}

	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	// Handle 401 — try refresh
	if resp.StatusCode == http.StatusUnauthorized && tokens.RefreshToken != "" {
		newTokens, refreshErr := auth.RefreshAccessToken(c.cfg.Server, c.cfg.ClientID, tokens)
		if refreshErr != nil {
			return fmt.Errorf("session expired: run 'sunshine auth login' to re-authenticate")
		}
		if err := auth.SaveTokens(newTokens); err != nil {
			return fmt.Errorf("saving refreshed tokens: %w", err)
		}

		// Retry with new token
		if bodyReader != nil {
			if body != nil {
				data, _ := json.Marshal(body)
				bodyReader = bytes.NewReader(data)
			}
		}
		req, _ = http.NewRequest(method, url, bodyReader)
		req.Header.Set("Authorization", "Bearer "+newTokens.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err = c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("retrying request: %w", err)
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp struct {
			Error   string            `json:"error"`
			Errors  map[string][]string `json:"errors"`
			Message string            `json:"message"`
		}
		if json.Unmarshal(respBody, &errResp) == nil {
			if errResp.Error != "" {
				return fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Error)
			}
			if errResp.Message != "" {
				return fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Message)
			}
		}
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}

	return nil
}

func (c *Client) getValidTokens() (*auth.TokenSet, error) {
	tokens, err := auth.LoadTokens()
	if err != nil || tokens == nil {
		return nil, fmt.Errorf("no tokens found")
	}

	if tokens.IsExpired() && tokens.RefreshToken != "" {
		newTokens, err := auth.RefreshAccessToken(c.cfg.Server, c.cfg.ClientID, tokens)
		if err != nil {
			return nil, err
		}
		if err := auth.SaveTokens(newTokens); err != nil {
			return nil, err
		}
		return newTokens, nil
	}

	return tokens, nil
}

// Scooters

func (c *Client) ListScooters(limit, offset int) ([]Scooter, error) {
	var scooters []Scooter
	path := "/scooters"
	params := []string{}
	if limit > 0 {
		params = append(params, fmt.Sprintf("limit=%d", limit))
	}
	if offset > 0 {
		params = append(params, fmt.Sprintf("offset=%d", offset))
	}
	if len(params) > 0 {
		path += "?" + strings.Join(params, "&")
	}
	err := c.get(path, &scooters)
	return scooters, err
}

func (c *Client) GetScooter(id int) (*Scooter, error) {
	var scooter Scooter
	err := c.get(fmt.Sprintf("/scooters/%d", id), &scooter)
	return &scooter, err
}

// Control commands

func (c *Client) Lock(id int) (*CommandResponse, error) {
	var resp CommandResponse
	err := c.post(fmt.Sprintf("/scooters/%d/lock", id), nil, &resp)
	return &resp, err
}

func (c *Client) Unlock(id int) (*CommandResponse, error) {
	var resp CommandResponse
	err := c.post(fmt.Sprintf("/scooters/%d/unlock", id), nil, &resp)
	return &resp, err
}

func (c *Client) Honk(id int) (*CommandResponse, error) {
	var resp CommandResponse
	err := c.post(fmt.Sprintf("/scooters/%d/honk", id), nil, &resp)
	return &resp, err
}

func (c *Client) SetBlinkers(id int, state string) (*CommandResponse, error) {
	var resp CommandResponse
	err := c.post(fmt.Sprintf("/scooters/%d/blinkers", id), map[string]string{"state": state}, &resp)
	return &resp, err
}

func (c *Client) OpenSeatbox(id int) (*CommandResponse, error) {
	var resp CommandResponse
	err := c.post(fmt.Sprintf("/scooters/%d/open_seatbox", id), nil, &resp)
	return &resp, err
}

func (c *Client) PlaySound(id int, sound string) (*CommandResponse, error) {
	var resp CommandResponse
	body := map[string]string{}
	if sound != "" {
		body["sound"] = sound
	}
	err := c.post(fmt.Sprintf("/scooters/%d/play_sound", id), body, &resp)
	return &resp, err
}

func (c *Client) Ping(id int) (*CommandResponse, error) {
	var resp CommandResponse
	err := c.post(fmt.Sprintf("/scooters/%d/ping", id), nil, &resp)
	return &resp, err
}

func (c *Client) GetState(id int) (*CommandResponse, error) {
	var resp CommandResponse
	err := c.post(fmt.Sprintf("/scooters/%d/get_state", id), nil, &resp)
	return &resp, err
}

func (c *Client) Locate(id int) (*CommandResponse, error) {
	var resp CommandResponse
	err := c.post(fmt.Sprintf("/scooters/%d/locate", id), nil, &resp)
	return &resp, err
}

func (c *Client) Alarm(id int, duration string) (*CommandResponse, error) {
	var resp CommandResponse
	body := map[string]string{}
	if duration != "" {
		body["duration"] = duration
	}
	err := c.post(fmt.Sprintf("/scooters/%d/alarm", id), body, &resp)
	return &resp, err
}

func (c *Client) AlarmArm(id int) (*CommandResponse, error) {
	var resp CommandResponse
	err := c.post(fmt.Sprintf("/scooters/%d/alarm_arm", id), nil, &resp)
	return &resp, err
}

func (c *Client) AlarmDisarm(id int) (*CommandResponse, error) {
	var resp CommandResponse
	err := c.post(fmt.Sprintf("/scooters/%d/alarm_disarm", id), nil, &resp)
	return &resp, err
}

func (c *Client) AlarmStop(id int) (*CommandResponse, error) {
	var resp CommandResponse
	err := c.post(fmt.Sprintf("/scooters/%d/alarm_stop", id), nil, &resp)
	return &resp, err
}

func (c *Client) Hibernate(id int) (*CommandResponse, error) {
	var resp CommandResponse
	err := c.post(fmt.Sprintf("/scooters/%d/hibernate", id), nil, &resp)
	return &resp, err
}

// Trips

func (c *Client) ListTrips(scooterID, limit, offset int) ([]Trip, error) {
	var trips []Trip
	path := fmt.Sprintf("/scooters/%d/trips", scooterID)
	params := []string{}
	if limit > 0 {
		params = append(params, fmt.Sprintf("limit=%d", limit))
	}
	if offset > 0 {
		params = append(params, fmt.Sprintf("offset=%d", offset))
	}
	if len(params) > 0 {
		path += "?" + strings.Join(params, "&")
	}
	err := c.get(path, &trips)
	return trips, err
}

func (c *Client) GetTrip(scooterID, tripID int) (*Trip, error) {
	var trip Trip
	err := c.get(fmt.Sprintf("/scooters/%d/trips/%d", scooterID, tripID), &trip)
	return &trip, err
}

// Destinations

func (c *Client) GetDestination(scooterID int) (*Destination, error) {
	var dest Destination
	err := c.get(fmt.Sprintf("/scooters/%d/destination", scooterID), &dest)
	return &dest, err
}

func (c *Client) SetDestination(scooterID int, lat, lng float64, address string) (*CommandResponse, error) {
	var resp CommandResponse
	body := map[string]interface{}{
		"latitude":  lat,
		"longitude": lng,
	}
	if address != "" {
		body["address"] = address
	}
	err := c.put(fmt.Sprintf("/scooters/%d/destination", scooterID), body, &resp)
	return &resp, err
}

func (c *Client) ClearDestination(scooterID int) error {
	return c.delete(fmt.Sprintf("/scooters/%d/destination", scooterID), nil)
}
