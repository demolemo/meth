package cmd

import (
	"fmt"

	"checking/store"

	"github.com/spf13/cobra"
)

var (
	metricCreateID   int
	metricCreateName string
)

var metricCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a metric",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		rowsAffected, err := store.CreateMetric(db, &store.Metric{ID: metricCreateID, Name: metricCreateName})
		if err != nil {
			return err
		}
		fmt.Printf("Metric created, rows affected: %d\n", rowsAffected)
		return nil
	},
}

func init() {
	metricCmd.AddCommand(metricCreateCmd)
	metricCreateCmd.Flags().IntVar(&metricCreateID, "id", 0, "metric ID")
	metricCreateCmd.Flags().StringVar(&metricCreateName, "name", "", "metric name")
	metricCreateCmd.MarkFlagRequired("id")
	metricCreateCmd.MarkFlagRequired("name")
}
