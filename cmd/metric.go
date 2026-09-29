package cmd

import (
	"fmt"

	"checking/store"

	"github.com/spf13/cobra"
)

var metricCmd = &cobra.Command{
	Use:   "metric",
	Short: "Manage metrics",
}

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
	rootCmd.AddCommand(metricCmd)

	metricCmd.AddCommand(metricCreateCmd)
	metricCreateCmd.Flags().IntVar(&metricCreateID, "id", 0, "metric ID")
	metricCreateCmd.Flags().StringVar(&metricCreateName, "name", "", "metric name")
	metricCreateCmd.MarkFlagRequired("id")
	metricCreateCmd.MarkFlagRequired("name")

	metricCmd.AddCommand(metricFindCmd)
	metricFindCmd.Flags().IntVar(&metricFindID, "id", 0, "metric ID")
	metricFindCmd.MarkFlagRequired("id")

	metricCmd.AddCommand(metricUpdateCmd)
	metricUpdateCmd.Flags().IntVar(&metricUpdateID, "id", 0, "metric ID")
	metricUpdateCmd.Flags().StringVar(&metricUpdateName, "name", "", "new metric name")
	metricUpdateCmd.MarkFlagRequired("id")
	metricUpdateCmd.MarkFlagRequired("name")

	metricCmd.AddCommand(metricDeleteCmd)
	metricDeleteCmd.Flags().IntVar(&metricDeleteID, "id", 0, "metric ID")
	metricDeleteCmd.MarkFlagRequired("id")

	metricCmd.AddCommand(metricListCmd)
}
