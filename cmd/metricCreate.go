package cmd

import (
	"fmt"
	"strconv"

	"checking/store"

	"github.com/spf13/cobra"
)

var metricCreateCmd = &cobra.Command{
	Use:   "create ID name",
	Short: "Create a metric",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		rowsAffected, err := store.CreateMetric(db, &store.Metric{ID: id, Name: args[1]})
		if err != nil {
			return err
		}
		fmt.Printf("Metric created, rows affected: %d\n", rowsAffected)
		return nil
	},
}

func init() {
	metricCmd.AddCommand(metricCreateCmd)
}
