package ioutil

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
)

// Get performs an HTTP GET request to the specified URL with optional query parameters,
// decodes the JSON response into a value of type T, and returns it.
//
// If client is nil, http.DefaultClient is used.
// If params is not nil and has values, they are appended to the URL as query string.
func Get[T any](client *http.Client, urlStr string, params url.Values) (T, error) {
	var result T

	if client == nil {
		client = http.DefaultClient
	}

	slog.Debug("fetching URL", "url", urlStr, "params", params)

	if len(params) > 0 {
		urlStr = urlStr + "?" + params.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return result, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}
