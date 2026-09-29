package cmd

import (
	"fmt"
	"strconv"

	"checking/store"

	"github.com/spf13/cobra"
)

var metricDeleteCmd = &cobra.Command{
	Use:   "delete ID",
	Short: "Delete a metric",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		rowsAffected, err := store.DeleteMetricByID(db, id)
		if err != nil {
			return err
		}
		fmt.Printf("Metric deleted, rows affected: %d\n", rowsAffected)
		return nil
	},
}

func init() {
	metricCmd.AddCommand(metricDeleteCmd)
}
