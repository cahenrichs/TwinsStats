package mlb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Client struct {
	baseURL string
	http    *http.Client
	cache   string
}

func NewClient() *Client {
	return &Client{
		baseURL: "https://statsapi.mlb.com/api/v1/",
		http:    &http.Client{},
		cache:   "cache",
	}
}

func fetchAndCache[T any](c *Client, url string, cachePath string) (*T, error) {
	fullCachePath := filepath.Join(c.cache, cachePath)

	if _, err := os.Stat(fullCachePath); err == nil {
		// Cache exists, return cached data
		data, err := os.ReadFile(fullCachePath)
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
	resp, err := c.http.Get(url)
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
	if err := os.MkdirAll(filepath.Dir(fullCachePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}
	if err := os.WriteFile(fullCachePath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write cache: %w", err)
	}
	return &result, nil
}

func (c *Client) FindPlayerByName(teamID int, playerName string) (*Person, error) {
	if playerName == "" {
		return nil, fmt.Errorf("player name cannot be empty")
	}
	roster, err := c.GetRoster(teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get roster: %w", err)
	}
	for _, entry := range roster.Roster {
		if strings.Contains(strings.ToLower(entry.Person.FullName), strings.ToLower(playerName)) {
			p := entry.Person
			return &p, nil
		}
	}
	return nil, fmt.Errorf("player not found: %s", playerName)
}

func (c *Client) GetRoster(teamID int) (*RosterResponse, error) {
	url := fmt.Sprintf("%steams/%d/roster", c.baseURL, teamID)
	cachePath := fmt.Sprintf("cache/roster_%d.json", teamID)
	return fetchAndCache[RosterResponse](c, url, cachePath)
}

func (c *Client) GetTeamStats(teamID int, season int) (*SeasonStatsResponse, error) {
	url := fmt.Sprintf("%steams/%d/stats?season=%d&group=hitting&stats=season", c.baseURL, teamID, season)
	cache := filepath.Join(c.cache, fmt.Sprintf("stats_hitting_%d_%d.json", teamID, season))
	return fetchAndCache[SeasonStatsResponse](c, url, cache)
}

func (c *Client) GetPitchingStats(playerID int) (*SeasonStatsResponse, error) {
	url := fmt.Sprintf("%speople/%d/stats?stats=yearByYear&group=pitching", c.baseURL, playerID)
	cachePath := fmt.Sprintf("cache/pitching_stats_%d.json", playerID)
	return fetchAndCache[SeasonStatsResponse](c, url, cachePath)
}

func (c *Client) GetPlayerSeasonStats(playerID int) (*SeasonStatsResponse, error) {
	url := fmt.Sprintf("%speople/%d/stats?stats=yearByYear", c.baseURL, playerID)
	cachePath := fmt.Sprintf("cache/player_stats_%d.json", playerID)
	return fetchAndCache[SeasonStatsResponse](c, url, cachePath)
}

func (c *Client) GetPlayerStatsbySeason(playerID int, season int) (*SeasonStatsResponse, error) {
	url := fmt.Sprintf("%speople/%d/stats?stats=season&season=%d&group=hitting", c.baseURL, playerID, season)
	cachePath := fmt.Sprintf("cache/player_stats_%d_%d.json", playerID, season)
	return fetchAndCache[SeasonStatsResponse](c, url, cachePath)
}
