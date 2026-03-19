package mlb

import (
	"encoding/json"
	"fmt"
)

// Get the full roster for a team, including player IDs and positions. This is used to get the player IDs for the next step of getting stats for each player.
type RosterResponse struct {
	Roster []RosterEntry `json:"roster"`
}

type RosterEntry struct {
	Person       Person `json:"person"`
	JerseyNumber string `json:"jerseyNumber"`
	Position     struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"position"`
}

//Get the stats for a player, including their hitting and pitching stats. This is used to get the stats for each player on the roster.

type SeasonStatsResponse struct {
	Stats []StatsContainer `json:"stats"`
}

type StatsContainer struct {
	Splits []StatSplit `json:"splits"`
}

type StatSplit struct {
	Season string          `json:"season"`
	Stat   json.RawMessage `json:"stat"`
	Player Person          `json:"player"`
}

type Person struct {
	ID              int    `json:"id"`
	FullName        string `json:"fullName"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	PrimaryPosition struct {
		Name string `json:"name"`
	} `json:"primaryPosition"`
}

type HittingStats struct {
	GamesPlayed    int    `json:"gamesPlayed"`
	AtBats         int    `json:"atBats"`
	Hits           int    `json:"hits"`
	BattingAverage string `json:"avg"`
	HomeRuns       int    `json:"homeRuns"`
	OBP            string `json:"obp"`
	SLG            string `json:"slg"`
	RBIs           int    `json:"rbi"`
	Runs           int    `json:"runs"`
	SB             int    `json:"stolenBases"`
}

type PitchingStats struct {
	GamesPlayed    int     `json:"gamesPlayed"`
	InningsPitched float64 `json:"inningsPitched"`
	Wins           int     `json:"wins"`
	Losses         int     `json:"losses"`
	ERA            float64 `json:"ERA"`
	WHIP           float64 `json:"WHIP"`
	Strikeouts     int     `json:"strikeouts"`
	SOP9           int     `json:"strikeoutsPer9Inn"`
}

func (s *StatSplit) GetHittingStats() (*HittingStats, error) {
	var h HittingStats
	err := json.Unmarshal(s.Stat, &h)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal hitting stats: %w", err)
	}
	return &h, nil
}

func (s *StatSplit) GetPitchingStats() (*PitchingStats, error) {
	var p PitchingStats
	err := json.Unmarshal(s.Stat, &p)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal pitching stats: %w", err)
	}
	return &p, nil
}
