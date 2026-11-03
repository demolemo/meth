package main

import (
	"database/sql"
	"os"

	"github.com/demolemo/meth"
	"github.com/demolemo/meth/config"
	"github.com/demolemo/meth/sqlite"
	"github.com/spf13/cobra"
)

// wired by PersistentPreRunE below, so they are ready before any command body
// runs and closed again after it
var (
	db *sql.DB
	ms meth.MetricService
	vs meth.ValueService
)

// --db, wins over the config file
var dbFlag string

var rootCmd = &cobra.Command{
	Use:   "meth",
	Short: "smol application for tracking METrics related to Health",
	Long:  `smol application for tracking METrics related to Health`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		path, err := dbPath()
		if err != nil {
			return err
		}
		// assignment, not :=, or these stay nil and every command panics
		db, err = sqlite.Open(path)
		if err != nil {
			return err
		}
		ms = sqlite.NewMetricService(db)
		vs = sqlite.NewValueService(db)
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if db == nil {
			return nil
		}
		return db.Close()
	},
}

// quietUsage stops cobra from printing the whole flag reference after a failure
// that has nothing to do with flags, like a metric that is not there.
//
// It has to be called inside RunE. Cobra validates required flags and flag
// groups *after* the PreRun hooks, so setting SilenceUsage any earlier, or on
// the root command, also swallows the usage block for genuine flag mistakes,
// which is the one case where printing it helps.
func quietUsage(cmd *cobra.Command) {
	cmd.SilenceUsage = true
}

// dbPath resolves the database location: --db first, then the config file,
// then the built-in default.
func dbPath() (string, error) {
	if dbFlag != "" {
		return dbFlag, nil
	}
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}
	return cfg.DB, nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dbFlag, "db", "", "sqlite file to use, overrides the config file")
}
