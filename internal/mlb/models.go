package mlb

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
	Splits StatSplit `json:"splits"`
}

type StatSplit struct {
	Stat   HittingStats `json:"stat"`
	Player Person       `json:"player"`
}

type Person struct {
	ID              int    `json:"id"`
	FullName        string `json:"fullName"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	PrimaryPosition string `json:"primaryPosition"`
}

type HittingStats struct {
	GamesPlayed    int     `json:"gamesPlayed"`
	AtBats         int     `json:"atBats"`
	Hits           int     `json:"hits"`
	BattingAverage float64 `json:"battingAverage"`
	HomeRuns       int     `json:"homeRuns"`
	OBP            float64 `json:"OBP"`
	SLG            float64 `json:"SLG"`
	RBIs           int     `json:"RBIs"`
	Runs           int     `json:"runs"`
	SB             int     `json:"stolenBases"`
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
