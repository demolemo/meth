package main

import (
	"fmt"
	"strings"

	"github.com/demolemo/meth"

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
		if strings.TrimSpace(metricCreateName) == "" {
			return fmt.Errorf("metric name can't be empty")
		}
		rowsAffected, err := ms.CreateMetric(&meth.Metric{ID: metricCreateID, Name: metricCreateName})
		if err != nil {
			return err
		}
		fmt.Printf("Metric created, rows affected: %d\n", rowsAffected)
		return nil
	},
}

var (
	metricFindID   int
	metricFindName string
)

var metricFindCmd = &cobra.Command{
	Use:   "find",
	Short: "Find a metric by ID",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		metricID, err := ms.ResolveMetricID(metricFindID, metricFindName)
		if err != nil {
			return err
		}

		met, err := ms.GetMetricByID(metricID)
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
		if strings.TrimSpace(metricUpdateName) == "" {
			return fmt.Errorf("metric name can't be empty")
		}
		rowsAffected, err := ms.UpdateMetricByID(&meth.MetricUpdate{ID: metricUpdateID, Name: metricUpdateName})
		if err != nil {
			return err
		}
		fmt.Printf("Metric updated, rows affected: %d\n", rowsAffected)
		return nil
	},
}

var (
	metricDeleteID   int
	metricDeleteName string
)

var metricDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a metric",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		metricID, err := ms.ResolveMetricID(metricDeleteID, metricDeleteName)
		if err != nil {
			return err
		}

		rowsAffected, err := ms.DeleteMetricByID(metricID)
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
		// NOTE: change for a function with limit, repeat once again
		mets, err := ms.GetMetrics()
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
	metricCreateCmd.Flags().StringVarP(&metricCreateName, "name", "n", "", "metric name")
	metricCreateCmd.MarkFlagRequired("id")
	metricCreateCmd.MarkFlagRequired("name")

	metricCmd.AddCommand(metricFindCmd)
	metricFindCmd.Flags().IntVar(&metricFindID, "id", 0, "metric ID")
	metricFindCmd.Flags().StringVarP(&metricFindName, "metric", "m", "", "metric name, alternative to id flag")
	metricFindCmd.MarkFlagsOneRequired("id", "metric")
	metricFindCmd.MarkFlagsMutuallyExclusive("id", "metric")

	metricCmd.AddCommand(metricUpdateCmd)
	metricUpdateCmd.Flags().IntVar(&metricUpdateID, "id", 0, "metric ID")
	metricUpdateCmd.Flags().StringVarP(&metricUpdateName, "name", "n", "", "new metric name")
	metricUpdateCmd.MarkFlagRequired("id")
	metricUpdateCmd.MarkFlagRequired("name")

	metricCmd.AddCommand(metricDeleteCmd)
	metricDeleteCmd.Flags().IntVar(&metricDeleteID, "id", 0, "metric ID")
	metricDeleteCmd.Flags().StringVarP(&metricDeleteName, "metric", "m", "", "metric name, alternative to id flag")
	metricDeleteCmd.MarkFlagsOneRequired("id", "metric")
	metricDeleteCmd.MarkFlagsMutuallyExclusive("id", "metric")

	metricCmd.AddCommand(metricListCmd)
}
