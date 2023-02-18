package hafas

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Stop struct {
	ExtID     string `json:"extId"`
	ID        string `json:"id"`
	ProdClass string `json:"prodClass"`
	State     string `json:"state"`
	Type      string `json:"type"`
	TypeStr   string `json:"typeStr"`
	Name      string `json:"value"`
	Weight    string `json:"weight"`
	Xcoord    string `json:"xcoord"`
	Ycoord    string `json:"ycoord"`
}

func parseStops(body io.Reader) ([]Stop, error) {
	bytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	stripped := strings.TrimSuffix(strings.TrimPrefix(string(bytes), "SLs.sls="), ";")

	parsed := struct {
		Suggestions []Stop `json:"suggestions"`
	}{}

	err = json.Unmarshal([]byte(stripped), &parsed)
	if err != nil {
		return nil, err
	}

	return parsed.Suggestions, nil
}

func (c *Client) GetStops(ctx context.Context, query string, nResults int) ([]Stop, error) {
	params := url.Values{}
	params.Add("encoding", "utf-8")
	params.Add("suggestMethod", "none")
	params.Add("S", query+"?")
	params.Add("REQ0JourneyStopsS0A", "1")
	params.Add("REQ0JourneyStopsB", fmt.Sprint(nResults))
	params.Add("REQ0JourneyStopsF", "distinguishStationAttribute;ZH")

	resp, err := http.Get(c.buildURL("/ajax-getstop.exe/dn", params))
	if err != nil {
		return nil, err
	}
	// #nosec G307
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, errors.New("non-200 status code")
	}

	return parseStops(resp.Body)
}

func (c *Client) FindStop(ctx context.Context, query string) (*Stop, error) {
	stops, err := c.GetStops(ctx, query, 1)
	if err != nil {
		return nil, err
	}

	if len(stops) < 1 {
		return nil, errors.New("no stop found")
	}

	return &stops[0], nil
}
