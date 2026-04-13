package gasbuddy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) GetCSRF(stationID string) (string, error) {
	url := c.baseURL + "/station/" + stationID

	req, _ := http.NewRequest("GET", url, nil)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	html := string(b)

	// simple extraction (we can harden later if needed)
	idx := strings.Index(html, "gbcsrf")
	if idx == -1 {
		return "", fmt.Errorf("gbcsrf not found")
	}

	frag := html[idx : idx+200]

	parts := strings.Split(frag, "\"")
	for _, p := range parts {
		if strings.HasPrefix(p, "1.") {
			return p, nil
		}
	}

	return "", fmt.Errorf("failed to parse gbcsrf")
}

func (c *Client) GetStationPrices(stationID string) ([]byte, error) {
	csrf, err := c.GetCSRF(stationID)
	if err != nil {
		return nil, err
	}

	query := `
	query GetStationPrices($id: ID!) {
		station(id: $id) {
			id
			prices {
				fuelProduct
				longName
				credit {
					price
					formattedPrice
					postedTime
				}
			}
		}
	}`

	payload := map[string]any{
		"operationName": "GetStationPrices",
		"variables": map[string]any{
			"id": stationID,
		},
		"query": query,
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", c.baseURL+"/graphql", bytes.NewBuffer(body))

	req.Header.Set("gbcsrf", csrf)

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bad status: %d body: %s", resp.StatusCode, string(b))
	}

	return io.ReadAll(resp.Body)
}
