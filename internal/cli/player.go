package cli

import (
	"fmt"

	"github.com/cahenrichs/mlbstats/internal/mlb"
	"github.com/cahenrichs/mlbstats/internal/models"
	"github.com/spf13/cobra"
)

const defaultYear = 2026

var playerYear int

var playerCmd = &cobra.Command{
	Use:   "Player",
	Short: "MLB Player Stats",
	Args:  cobra.ExactArgs(1),
	RunE:  runPlayerCmd,
}

func init() {
	playerCmd.Flags().IntVar(&playerYear, "year", defaultYear, "Season year")
}

func runPlayerCmd(cmd *cobra.Command, args []string) error {
	playerName := args[0]
	repo, err := initDB()
	if err != nil {
		return err
	}
	api := mlb.NewClient()

	if verbose {
		fmt.Printf("searching for player: %s\n", playerName)
	}

	player, err := repo.FindPlayerByName(playerName)
	if err != nil {
		if verbose {
			fmt.Println("Player not found in DB, searching MLB API")
		}
		found, err := api.FindPlayerAllTeams(playerName)
		if err != nil {
			return fmt.Errorf("player not found: %s", playerName)
		}
		if len(found) == 0 {
			return fmt.Errorf("player not found: %s", playerName)
		}
		entry := found[0]
		player = &models.Player{
			MLBID:    entry.Person.ID,
			FullName: entry.Person.FullName,
			Position: entry.Person.PrimaryPosition.Name,
		}
		if err := repo.SavePlayer(player); err != nil {
			return fmt.Errorf("failed to save player: %w", err)
		}
		if verbose {
			fmt.Printf("Saved player to DB: %s\n", player.FullName)
		}
	}
	stats, err := repo.GetPlayerStats(int(player.Id), playerYear, player.Position)
	if err != nil {
		if verbose {
			fmt.Println("Stats not inDB, fetching from MLB")
		}

		group := "hitting"
		if player.IsPitcher() {
			group = "pitching"
		}
		apiStats, err := api.GetPlayerStats(player.MLBID, playerYear, group)
		if err != nil {
			return fmt.Errorf("Failed to get stats: %w\n", err)
		}
		stats = apiStats
	}
	fmt.Printf("\n=== %s (%d) ===\n", player.FullName, player.MLBID)
	fmt.Printf("Position: %s\n\n", player.Position)

	if player.IsPitcher() {
		printPitchingStats(stats.(*mlb.SeasonStatsResponse), playerYear)
	} else {
		printHittingStats(stats.(*mlb.SeasonStatsResponse), playerYear)
	}
	return nil
}

func printPitchingStats(stats *mlb.SeasonStatsResponse, year int) {
	for _, container := range stats.Stats {
		for _, split := range container.Splits {
			if split.Season != fmt.Sprintf("%d", year) {
				continue
			}
			p, err := split.GetPitchingStats()
			if err != nil {
				fmt.Printf("Error parsing stats: %v\n", err)
				continue
			}
			fmt.Printf("Pitching Stats (%s)\n", split.Season)
			fmt.Printf("  Games:   %d\n", p.GamesPlayed)
			fmt.Printf("  IP:      %s\n", p.InningsPitched)
			fmt.Printf("  W-L:     %d-%d\n", p.Wins, p.Losses)
			fmt.Printf("  ERA:     %s\n", p.ERA)
			fmt.Printf("  WHIP:    %s\n", p.WHIP)
			fmt.Printf("  K:       %d\n", p.Strikeouts)
			fmt.Printf("  BB:      %d\n", p.Walks)
			fmt.Printf("  H:       %d\n", p.Hits)
			fmt.Printf("  R:       %d\n", p.Runs)
			fmt.Printf("  HR:      %d\n", p.HomeRuns)
			fmt.Printf("  SV:      %d\n", p.Saves)
			fmt.Printf("  K/9:     %s\n", p.SOP9)
			return
		}
	}
	fmt.Printf("No pitching stats found for %d\n", year)
}

func printHittingStats(stats *mlb.SeasonStatsResponse, year int) {
	for _, container := range stats.Stats {
		for _, split := range container.Splits {
			if split.Season != fmt.Sprintf("%d", year) {
				continue
			}
			h, err := split.GetHittingStats()
			if err != nil {
				fmt.Printf("Error parsing stats: %v\n", err)
				continue
			}
			fmt.Printf("Hitting Stats (%s)\n", split.Season)
			fmt.Printf("  Games:   %d\n", h.GamesPlayed)
			fmt.Printf("  At Bats: %d\n", h.AtBats)
			fmt.Printf("  Hits:    %d\n", h.Hits)
			fmt.Printf("  BA:      %s\n", h.BattingAverage)
			fmt.Printf("  HR:      %d\n", h.HomeRuns)
			fmt.Printf("  RBI:     %d\n", h.RBIs)
			fmt.Printf("  Runs:    %d\n", h.Runs)
			fmt.Printf("  OBP:     %s\n", h.OBP)
			fmt.Printf("  SLG:     %s\n", h.SLG)
			return
		}
	}
	fmt.Printf("No hitting stats found for %d\n", year)
}
