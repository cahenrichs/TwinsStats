package cli

import (
	"fmt"
	"strconv"

	"github.com/cahenrichs/mlbstats/internal/mlb"
	"github.com/cahenrichs/mlbstats/internal/models"
	"github.com/spf13/cobra"
)

var teamYear int

var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Mlb team stats",
}

var statsCmd = &cobra.Command{
	Use:   "stats <team-name>",
	Short: "Display team stats for a season",
	Args:  cobra.ExactArgs(1),
	RunE:  runTeamStats,
}

func init() {
	teamCmd.AddCommand(statsCmd)
	statsCmd.Flags().IntVar(&teamYear, "year", defaultYear, "Season year")
}

func runTeamStats(cmd *cobra.Command, args []string) error {
	teamName := args[0]
	repo, err := initDB()
	if err != nil {
		return err
	}
	api := mlb.NewClient()

	if verbose {
		fmt.Printf("Searching for team: %s\n", teamName)
	}

	team, err := repo.FindTeamByName(teamName)
	if err != nil {
		if verbose {
			fmt.Println("Team not found in DB, searching MLB...")
		}
		matches, err := api.SearchTeams(teamName)
		if err != nil {
			return fmt.Errorf("failed to search teams: %w", err)
		}
		if len(matches) == 0 {
			return fmt.Errorf("team not found: %s", teamName)
		}
		m := matches[0]
		team = &models.Team{
			MLBID: m.ID,
			Name:  m.Name,
		}
		if err := repo.SaveTeam(team); err != nil {
			return fmt.Errorf("failed to save team: %w", err)
		}
		if verbose {
			fmt.Printf("Saved team to DB: %s\n", team.Name)
		}
	}
	if verbose {
		fmt.Printf("Fetching stats for %s", team.Name)
	}
	fmt.Printf("\n=== %s (%d) ===\n\n", team.Name, team.MLBID)

	hitting, err := api.GetTeamStats(team.MLBID, teamYear, "hitting")
	if err != nil {
		return fmt.Errorf("failed to get team hitting stats: %w", err)
	}
	printTeamHittingStats(hitting, teamYear)

	pitching, err := api.GetTeamStats(team.MLBID, teamYear, "pitching")
	if err != nil {
		return fmt.Errorf("failed to get team pitching stats: %w", err)
	}
	printTeamPitchingStats(pitching, teamYear)
	return nil
}

func printTeamHittingStats(stats *mlb.SeasonStatsResponse, year int) {
	for _, container := range stats.Stats {
		for _, split := range container.Splits {
			if split.Season != strconv.Itoa(year) {
				continue
			}
			h, err := split.GetHittingStats()
			if err != nil {
				fmt.Printf("Error parsing hitting stats: %v\n", err)
				return
			}
			fmt.Printf("Team Hitting (%s)\n", split.Season)
			fmt.Printf("  Games:   %d\n", h.GamesPlayed)
			fmt.Printf("  At Bats: %d\n", h.AtBats)
			fmt.Printf("  Hits:    %d\n", h.Hits)
			fmt.Printf("  BA:      %s\n", h.BattingAverage)
			fmt.Printf("  HR:      %d\n", h.HomeRuns)
			fmt.Printf("  RBI:     %d\n", h.RBIs)
			fmt.Printf("  Runs:    %d\n", h.Runs)
			fmt.Printf("  OBP:     %s\n", h.OBP)
			fmt.Printf("  SLG:     %s\n", h.SLG)
			fmt.Printf("  SB:      %d\n\n", h.SB)
			return
		}
	}
	fmt.Printf("No team hiitting stats found for %d\n\n", year)
}

func printTeamPitchingStats(stats *mlb.SeasonStatsResponse, year int) {
	for _, container := range stats.Stats {
		for _, split := range container.Splits {
			if split.Season != strconv.Itoa(year) {
				continue
			}
			p, err := split.GetPitchingStats()
			if err != nil {
				fmt.Printf("Error parsing pitching stats: %v\n", err)
				return
			}
			fmt.Printf("Team Pitching (%s)\n", split.Season)
			fmt.Printf("  Games:   %d\n", p.GamesPlayed)
			fmt.Printf("  GS:      %d\n", p.GamesStarted)
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
	fmt.Printf("No team pitching stats found for %d\n", year)
}
