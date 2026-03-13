package mlb

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
