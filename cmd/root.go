/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"database/sql"
	"os"

	"github.com/demolemo/meth/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

// opened before any command runs, see PersistentPreRunE
var ms *sqlite.MetricService
var vs *sqlite.ValueService
var db *sql.DB

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "meth",
	Short: "smol application for tracking METrics related to Health",
	Long:  `smol application for tracking METrics related to Health`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// db should be wired here temporarily just for testing
		db, err := sql.Open("sqlite3", "test.db?_foreign_keys=on")
		if err != nil {
			return err
		}
		// metric service is wired for all commands
		ms = sqlite.NewMetricService(db)
		vs = sqlite.NewValueService(db)
		return err
	},
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
	// NOTE: remove all that bullshit. that is legacy code from cobra
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.meth.yaml)")
}
