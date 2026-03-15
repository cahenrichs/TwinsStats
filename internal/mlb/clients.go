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

func GetRoster(teamID int) (*RosterResponse, error) {
	url := fmt.Sprintf("%steams/%d/roster", BaseURL, teamID)
	cachePath := fmt.Sprintf("cache/roster_%d.json", teamID)
	return fetchAndCache[RosterResponse](url, cachePath)
}

func GetTeamStats(teamID int, season int) (*SeasonStatsResponse, error) {
	url := fmt.Sprintf("%steams/%d/stats?season=%d&group=hitting&stats=season", BaseURL, teamID, season)
	cache := filepath.Join("data", fmt.Sprintf("stats_hitting_%d_%d.json", teamID, season))
	return fetchAndCache[SeasonStatsResponse](url, cache)
}

func GetPitchingStats(playerID int) (*SeasonStatsResponse, error) {
	url := fmt.Sprintf("%speople/%d/stats?stats=yearByYear&group=pitching", BaseURL, playerID)
	cachePath := fmt.Sprintf("cache/pitching_stats_%d.json", playerID)
	return fetchAndCache[SeasonStatsResponse](url, cachePath)
}

func GetPlayerSeasonStats(playerID int) (*SeasonStatsResponse, error) {
	url := fmt.Sprintf("%speople/%d/stats?stats=yearByYear", BaseURL, playerID)
	cachePath := fmt.Sprintf("cache/player_stats_%d.json", playerID)
	return fetchAndCache[SeasonStatsResponse](url, cachePath)
}
func GetPlayerStatsbySeason(playerID int, season int) (*SeasonStatsResponse, error) {
	url := fmt.Sprintf("%speople/%d/stats?stats=season&season=%d&group=hitting", BaseURL, playerID, season)
	cachePath := fmt.Sprintf("cache/player_stats_%d_%d.json", playerID, season)
	return fetchAndCache[SeasonStatsResponse](url, cachePath)
}
