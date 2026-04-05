package cli

import (
	"fmt"
	"os"

	"github.com/cahenrichs/TwinsStats/internal/db"
	"github.com/spf13/cobra"
)

var (
	dbPath  string
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "mlb",
	Short: "Mlb Stats CLI - Query team and player statistics",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", "./data/mlb.db", "path to SQLite database")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func initDB() (*db.Repository, error) {
	gormDB, err := db.InitDB(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initalize database: %w", err)
	}
	return db.NewRepository(gormDB), nil
}
