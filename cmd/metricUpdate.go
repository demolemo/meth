package cmd

import (
	"fmt"
	"strconv"

	"checking/store"

	"github.com/spf13/cobra"
)

var metricUpdateCmd = &cobra.Command{
	Use:   "update ID new-name",
	Short: "Rename a metric",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		rowsAffected, err := store.UpdateMetricByID(db, &store.MetricUpdate{ID: id, Name: args[1]})
		if err != nil {
			return err
		}
		fmt.Printf("Metric updated, rows affected: %d\n", rowsAffected)
		return nil
	},
}

func init() {
	metricCmd.AddCommand(metricUpdateCmd)
}
