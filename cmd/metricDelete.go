package cmd

import (
	"fmt"

	"checking/store"

	"github.com/spf13/cobra"
)

var metricDeleteID int

var metricDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a metric",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		rowsAffected, err := store.DeleteMetricByID(db, metricDeleteID)
		if err != nil {
			return err
		}
		fmt.Printf("Metric deleted, rows affected: %d\n", rowsAffected)
		return nil
	},
}

func init() {
	metricCmd.AddCommand(metricDeleteCmd)
	metricDeleteCmd.Flags().IntVar(&metricDeleteID, "id", 0, "metric ID")
	metricDeleteCmd.MarkFlagRequired("id")
}
