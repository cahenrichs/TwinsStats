package main

import (
	"fmt"

	"github.com/cahenrichs/TwinsStats/internal/mlb"
)

func main() {

	fmt.Println("Fetching Minnesota Twins Stats...")

	api := mlb.Client{}

	stats, err := api.GetTeamStats(142, 2025)
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

	//get Buxton's stats
	fmt.Println("\nFetching Buxton's Stats...")
	searchName := "Buxton"
	player, err := api.FindPlayerByName(142, searchName)
	if err != nil {
		fmt.Printf("Error finding player: %v\n", err)
		return
	}
	fmt.Printf("Found player: %s (ID: %d)\n", player.FullName, player.ID)

	// Get Buxton's season stats for 2025
	playerStats, err := api.GetPlayerStatsbySeason(player.ID, 2025)
	if err != nil {
		fmt.Printf("Error fetching player stats: %v\n", err)
		return
	}
	if len(playerStats.Stats) > 0 && len(playerStats.Stats[0].Splits) > 0 {
		stat := playerStats.Stats[0].Splits[0].Stat
		fmt.Printf("Season: %s, Games Played: %d, At Bats: %d, Hits: %d, Batting Average: %s, Home Runs: %d\n",
			playerStats.Stats[0].Splits[0].Season, stat.GamesPlayed, stat.AtBats, stat.Hits, stat.BattingAverage, stat.HomeRuns)
	} else {
		fmt.Println("No stats found for Buxton in 2025.")
	}
}
