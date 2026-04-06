package cli

import "github.com/spf13/cobra"

const defaultYear = 2026

var playerCmd = &cobra.Command{
	Use:   "Player",
	Short: "MLB Player Stats",
}
