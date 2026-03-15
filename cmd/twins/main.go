package main

import (
	"fmt"

	"github.com/cahenrichs/TwinsStats/internal/mlb"
)

func main() {

	fmt.Println("Fetching Minnesota Twins Stats...")

	stats, err := mlb.GetTeamStats(142, 2025)
	if err != nil {
		fmt.Printf("Error fetching team stats: %v\n", err)
		return
	}
	for _, container := range stats.Stats {
		for _, stat := range container.Splits {
			fmt.Printf("Season: %s, Games Played: %d, At Bats: %d, Hits: %d, Batting Average: %s\n",
				stat.Season, stat.Stat.GamesPlayed, stat.Stat.AtBats, stat.Stat.Hits, stat.Stat.BattingAverage)
		}
	}
}
