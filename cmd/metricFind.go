package cmd

import (
	"fmt"
	"strconv"

	"checking/store"

	"github.com/spf13/cobra"
)

var metricFindCmd = &cobra.Command{
	Use:   "find ID",
	Short: "Find a metric by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		met, err := store.GetMetricByID(db, id)
		if err != nil {
			return err
		}
		if met == nil {
			return fmt.Errorf("metric %d not found", id)
		}
		fmt.Printf("Metric, ID: %d, Name: %s\n", met.ID, met.Name)
		return nil
	},
}

func init() {
	metricCmd.AddCommand(metricFindCmd)
}
