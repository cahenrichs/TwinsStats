package models

import "time"

type Team struct {
	Id        uint   `gorm:"primarykey"`
	MLBID     int    `gorm:"uniqueindex"`
	Name      string `gorm:"index"`
	Nickname  string `gorm:"index"`
	Abbr      string `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Player struct {
	Id            uint   `gorm:"primarykey"`
	MLBID         int    `gorm:"uniqueindex"`
	FullName      string `gorm:"index"`
	Position      string
	PositionCode  string
	CurrentTeamId uint            `gorm:"index"`
	HittingStats  []HittingStats  `gorm:"foreignKey:PlayerId"`
	PitchingStats []PitchingStats `gorm:"foreignKey:PlayerId"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type HittingStats struct {
	ID          uint `gorm:"primaryKey"`
	PlayerId    uint `gorm:"index"`
	Season      int  `gorm:"index"`
	GamesPlayed int
	AtBats      int
	Hits        int
	Doubles     int
	Triples     int
	HR          int
	RBI         int
	Runs        int
	SB          int
	Walks       int
	Strikeouts  int
	BA          float64
	OBP         float64
	SLG         float64
	OPS         float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PitchingStats struct {
	ID              uint `gorm:"primaryKey"`
	PlayerId        uint `gorm:"index"`
	Season          int  `gorm:"index"`
	GamesPlayed     int
	GamesStarted    int
	InningsPitched  float64
	HitsAllowed     int
	WalksAllowed    int
	RunsAllowed     int
	HomeRunsAllowed int
	Wins            int
	Losses          int
	Saves           int
	ERA             float64
	WHIP            float64
	Strikeouts      int
	SOP9            float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (p *Player) IsPitcher() bool {
	return p.PositionCode == "1"
}
