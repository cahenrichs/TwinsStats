package models

type Team struct {
	Id       int    `gorm:"primarykey"`
	Name     string `gorm:"index"`
	Nickname string `gorm:"index"`
	Abbr     string `gorm:"index"`
}

type Player struct {
	Id            int    `gorm:"primarykey"`
	FullName      string `gorm:"index"`
	Position      string
	CurrentTeamId int  `gorm:"index"`
	Team          Team `gorm:"foreignKey:CurrentTeamId"`
	HittingStats  []HittingStats
	PitchingStats []PitchingStats
}

type HittingStats struct {
	PlayerId    int  `gorm:"primarykey"`
	Season      int  `gorm:"primarykey"`
	IsCareer    bool `gorm:"primarykey"`
	GamesPlayed int
	AtBats      int
	Hits        int
	Doubles     int
	Triples     int
	HR          int
	RBI         int
	Runs        int
	SB          int
	BA          float64
	OBP         float64
	SLG         float64
	OPS         float64
}

type PitchingStats struct {
	PlayerId       int  `gorm:"primarykey"`
	Season         int  `gorm:"primarykey"`
	IsCareer       bool `gorm:"primarykey"`
	GamesPlayed    int
	GamesStarted   int
	InningsPitched float64
	Wins           int
	Losses         int
	Saves          int
	ERA            float64
	WHIP           float64
	Strikeouts     int
	SOP9           int
}
