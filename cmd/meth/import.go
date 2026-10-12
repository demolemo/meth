package main

import (
	"fmt"

	"github.com/demolemo/meth/parse"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a dump that brings its own metric names",
	Long: `Import a dump that brings its own metric names.

Unlike "value import", which files a file of numbers under one metric you
name, these sources decide the names themselves and create what they need.`,
}

var importWhoopFile string

var importWhoopCmd = &cobra.Command{
	Use:   "whoop",
	Short: "Import a whoop dump into the sleep., nap. and recovery. families",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// parse it all before touching the db - a bad row imports nothing
		readings, err := parse.Whoop(importWhoopFile)
		if err != nil {
			return err
		}

		metrics := 0
		mvs, err := parse.Values(readings, func(name, unit string) (int, error) {
			m, err := ms.EnsureMetric(name, unit)
			if err != nil {
				return 0, err
			}
			metrics++
			return m.ID, nil
		})
		if err != nil {
			return err
		}

		res, err := vs.CreateMetricValues(mvs)
		if err != nil {
			return err
		}
		fmt.Printf("Imported %s: %d metrics, %d readings -> %d new, %d updated, %d already there\n",
			importWhoopFile, metrics, len(mvs), res.Inserted, res.Updated, res.Unchanged)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.AddCommand(importWhoopCmd)
	importWhoopCmd.Flags().StringVarP(&importWhoopFile, "file", "f", "", "whoop dump .json")
	importWhoopCmd.MarkFlagRequired("file")
}
