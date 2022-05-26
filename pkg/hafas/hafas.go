package hafas

import (
	"errors"
	"net/url"
	"time"
)

/*
Package hafas implements a Go client for a (small)
subset of the HAFAS REST API.

Use NewClient() to get started.
*/

// Client is an API client for a HAFAS REST API endpoint.
type Client struct {
	// those are immutable after instantiation
	baseURL  url.URL
	timeZone *time.Location
}

// NewClient creates a new client.
// Parameters:
//   baseURL: e.g. https://zvv.hafas.cloud/bin
//   timeZone: e.g. Europe/Zurich
func NewClient(baseURL, timeZone string) (*Client, error) {

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if u.Path != "bin" {
		return nil, errors.New("invalid base URL")
	}

	tz, err := time.LoadLocation(timeZone)
	if err != nil {
		return nil, err
	}

	return &Client{baseURL: *u, timeZone: tz}, nil
}

func (c *Client) buildURL(path string, params url.Values) string {
	u := c.baseURL

	u.Path += path
	u.RawQuery = params.Encode()

	return u.String()
}
