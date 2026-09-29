package cmd

import (
	"fmt"

	"checking/store"

	"github.com/spf13/cobra"
)

var (
	metricUpdateID   int
	metricUpdateName string
)

var metricUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Rename a metric",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		rowsAffected, err := store.UpdateMetricByID(db, &store.MetricUpdate{ID: metricUpdateID, Name: metricUpdateName})
		if err != nil {
			return err
		}
		fmt.Printf("Metric updated, rows affected: %d\n", rowsAffected)
		return nil
	},
}

func init() {
	metricCmd.AddCommand(metricUpdateCmd)
	metricUpdateCmd.Flags().IntVar(&metricUpdateID, "id", 0, "metric ID")
	metricUpdateCmd.Flags().StringVar(&metricUpdateName, "name", "", "new metric name")
	metricUpdateCmd.MarkFlagRequired("id")
	metricUpdateCmd.MarkFlagRequired("name")
}
