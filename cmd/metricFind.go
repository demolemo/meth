package cmd

import (
	"fmt"

	"checking/store"

	"github.com/spf13/cobra"
)

var metricFindID int

var metricFindCmd = &cobra.Command{
	Use:   "find",
	Short: "Find a metric by ID",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		met, err := store.GetMetricByID(db, metricFindID)
		if err != nil {
			return err
		}
		if met == nil {
			return fmt.Errorf("metric %d not found", metricFindID)
		}
		fmt.Printf("Metric, ID: %d, Name: %s\n", met.ID, met.Name)
		return nil
	},
}

func init() {
	metricCmd.AddCommand(metricFindCmd)
	metricFindCmd.Flags().IntVar(&metricFindID, "id", 0, "metric ID")
	metricFindCmd.MarkFlagRequired("id")
}
