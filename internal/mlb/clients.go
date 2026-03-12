package mlb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

const BaseURL = "https://statsapi.mlb.com/api/v1/"

func fetchAndCache[T any](BaseURL string, cachePath string) (*T, error) {
	if _, err := os.Stat(cachePath); err == nil {
		// Cache exists, return cached data
		data, err := os.ReadFile(cachePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read cache: %w", err)
		}
		var result T
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cache: %w", err)
		}
		return &result, nil
	}
	//If not stored get it from the API
	resp, err := http.Get(BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Cache the result for future use and make sure dir exists
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write cache: %w", err)
	}
	return &result, nil
}
