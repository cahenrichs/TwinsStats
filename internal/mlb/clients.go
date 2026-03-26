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

func fetchAndCache[T any](c *Client, url string, fileName string) (*T, error) {
	fullCachePath := filepath.Join(c.cache, fileName)

	/*if _, err := os.Stat(fullCachePath); err == nil {
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
	} */
	if data, err := os.ReadFile(fullCachePath); err == nil {
		var result T
		if err := json.Unmarshal(data, &result); err == nil {
			return &result, nil
		}
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

	// Async or background caching could be done here, but for simplicity:
	_ = os.MkdirAll(filepath.Dir(fullCachePath), 0755)
	if data, err := json.MarshalIndent(result, "", "  "); err == nil {
		_ = os.WriteFile(fullCachePath, data, 0644)
	}
	return &result, nil

}

// Cache the result for future use and make sure dir exists
/*if err := os.MkdirAll(filepath.Dir(fullCachePath), 0755); err != nil {
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
}*/

func (c *Client) FindPlayerByName(teamID int, playerName string) (*RosterEntry, error) {
	if playerName == "" {
		return nil, fmt.Errorf("player name cannot be empty")
	}
	roster, err := c.GetRoster(teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get roster: %w", err)
	}
	for _, entry := range roster.Roster {
		if strings.Contains(strings.ToLower(entry.Person.FullName), strings.ToLower(playerName)) {
			return &entry, nil
		}
	}
	return nil, fmt.Errorf("player not found: %s", playerName)
}

func (c *Client) GetRoster(teamID int) (*RosterResponse, error) {
	url := fmt.Sprintf("%steams/%d/roster?hydrate=person(firstName,lastName,primaryPosition)", c.baseURL, teamID)
	fileName := fmt.Sprintf("roster_%d.json", teamID)
	return fetchAndCache[RosterResponse](c, url, fileName)
}

func (c *Client) GetTeamStats(teamID int, season int) (*SeasonStatsResponse, error) {
	url := fmt.Sprintf("%steams/%d/stats?season=%d&group=hitting&stats=season", c.baseURL, teamID, season)
	fileName := fmt.Sprintf("team_stats_%d_%d.json", teamID, season)
	return fetchAndCache[SeasonStatsResponse](c, url, fileName)
}

func (c *Client) GetPitchingStats(playerID int) (*SeasonStatsResponse, error) {
	url := fmt.Sprintf("%speople/%d/stats?stats=yearByYear&group=pitching", c.baseURL, playerID)
	fileName := fmt.Sprintf("pitching_stats_%d.json", playerID)
	return fetchAndCache[SeasonStatsResponse](c, url, fileName)
}

func (c *Client) GetPlayerSeasonStats(playerID int) (*SeasonStatsResponse, error) {
	url := fmt.Sprintf("%speople/%d/stats?stats=yearByYear", c.baseURL, playerID)
	fileName := fmt.Sprintf("player_stats_%d.json", playerID)
	return fetchAndCache[SeasonStatsResponse](c, url, fileName)
}

/*func (c *Client) GetPlayerStatsbySeason(playerID int, season int) (*SeasonStatsResponse, error) {
	url := fmt.Sprintf("%speople/%d/stats?stats=season&season=%d&group=hitting", c.baseURL, playerID, season)
	cachePath := fmt.Sprintf("cache/player_stats_%d_%d.json", playerID, season)
	return fetchAndCache[SeasonStatsResponse](c, url, cachePath)
} */

func (c *Client) GetPlayerStats(playerID int, season int, group string) (*SeasonStatsResponse, error) {
	url := fmt.Sprintf("%speople/%d/stats?stats=season&season=%d&group=%s&hydrate=person(firstName,lastName,primaryPosition)", c.baseURL, playerID, season, group)
	fileName := fmt.Sprintf("player_stats_%d_%d_%s.json", playerID, season, group)
	return fetchAndCache[SeasonStatsResponse](c, url, fileName)
}
