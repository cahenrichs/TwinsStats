package main

import (
	"fmt"

	"github.com/cahenrichs/TwinsStats/internal/mlb"
)

func main() {

	fmt.Println("Fetching Minnesota Twins Stats...")

	api := mlb.NewClient()

	stats, err := api.GetTeamStats(142, 2025)
	if err != nil {
		fmt.Printf("Error fetching team stats: %v\n", err)
		return
	}
	for _, container := range stats.Stats {
		for _, stat := range container.Splits {
			hStats, err := stat.GetHittingStats()
			if err != nil {
				fmt.Printf("Error parsing hitting stats: %v\n", err)
				continue
			}
			fmt.Printf("Season: %s, Games Played: %d, At Bats: %d, Hits: %d, Batting Average: %s\n",
				stat.Season, hStats.GamesPlayed, hStats.AtBats, hStats.Hits, hStats.BattingAverage)
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
	fmt.Printf("Found player: %s (ID: %d)\n", player.Person.FullName, player.Person.ID)

	// Get Buxton's season stats for 2025
	playerStats, err := api.GetPlayerStats(player.Person.ID, 2025, "hitting")
	if err != nil {
		fmt.Printf("Error fetching player stats: %v\n", err)
		return
	}
	if len(playerStats.Stats) > 0 && len(playerStats.Stats[0].Splits) > 0 {
		split := playerStats.Stats[0].Splits[0]
		hStat, err := split.GetHittingStats()
		if err != nil {
			fmt.Printf("Error parsing hitting stats: %v\n", err)
			return
		}
		fmt.Printf("Season: %s, Games Played: %d, At Bats: %d, Hits: %d, Batting Average: %s, Home Runs: %d\n",
			playerStats.Stats[0].Splits[0].Season, hStat.GamesPlayed, hStat.AtBats, hStat.Hits, hStat.BattingAverage, hStat.HomeRuns)
	} else {
		fmt.Println("No stats found for Buxton in 2025.")
	}
}
