package cmd

import (
	"fmt"

	"checking/store"

	"github.com/spf13/cobra"
)

var metricListCmd = &cobra.Command{
	Use:   "list",
	Short: "List metrics (first 5)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		mets, err := store.GetMetrics(db)
		if err != nil {
			return err
		}
		if len(mets) > 5 {
			fmt.Printf("Following metrics exist (cut to first 5):\n")
			mets = mets[:5]
		} else {
			fmt.Printf("Following metrics exist:\n")
		}
		for _, m := range mets {
			fmt.Printf("\tMetric, ID: %d, Name: %s\n", m.ID, m.Name)
		}
		return nil
	},
}

func init() {
	metricCmd.AddCommand(metricListCmd)
}
