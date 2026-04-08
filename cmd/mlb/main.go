package main

import (
	"github.com/cahenrichs/mlbstats/internal/cli"
)

func main() {
	cli.Execute()
	/*if len(os.Args) < 2 {
		fmt.Println("Usage: mlb <team-name>")
		fmt.Println("Example: mlb Dodgers")
		os.Exit(1)
	}

	teamName := os.Args[1]
	api := mlb.NewClient()

	teams, err := api.SearchTeams(teamName)
	if err != nil {
		fmt.Printf("Error searching teams: %v\n", err)
		os.Exit(1)
	}

	if len(teams) == 0 {
		fmt.Printf("No team found matching: %s\n", teamName)
		os.Exit(1)
	}

	team := teams[0]
	fmt.Printf("Fetching %s Stats...\n\n", team.Name)

	stats, err := api.GetTeamStats(team.ID, 2025, "hitting")
	if err != nil {
		fmt.Printf("Error fetching team stats: %v\n", err)
		os.Exit(1)
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
	} */
}
