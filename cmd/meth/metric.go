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
	metricCreateName string
)

var metricCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a metric",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		quietUsage(cmd)
		if strings.TrimSpace(metricCreateName) == "" {
			return fmt.Errorf("metric name can't be empty")
		}
		err := ms.CreateMetric(&meth.Metric{Name: metricCreateName})
		if err != nil {
			return err
		}
		fmt.Printf("Metric created\n")
		return nil
	},
}

var metricFindName string

var metricFindCmd = &cobra.Command{
	Use:   "find",
	Short: "Find a metric by ID",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		quietUsage(cmd)
		met, err := ms.GetMetric(metricFindName)
		if err != nil {
			return err
		}
		fmt.Printf("Metric, ID: %d, Name: %s\n", met.ID, met.Name)
		return nil
	},
}

var metricDeleteName string

var metricDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a metric",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		quietUsage(cmd)
		err := ms.DeleteMetric(metricDeleteName)
		if err != nil {
			return err
		}
		fmt.Printf("Metric deleted\n")
		return nil
	},
}

var metricListLimit int

var metricListCmd = &cobra.Command{
	Use:   "list",
	Short: "List metrics (first 5)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		quietUsage(cmd)
		mets, err := ms.GetMetrics(metricListLimit)
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
	// NOTE: a little bit poor job because we name NAME METRIC but this keeps everything else consisstent
	metricCreateCmd.Flags().StringVarP(&metricCreateName, "metric", "m", "", "metric name")
	metricCreateCmd.MarkFlagRequired("name")

	metricCmd.AddCommand(metricFindCmd)
	metricFindCmd.Flags().StringVarP(&metricFindName, "metric", "m", "", "metric name, alternative to id flag")
	metricFindCmd.MarkFlagRequired("metric")

	metricCmd.AddCommand(metricDeleteCmd)
	metricDeleteCmd.Flags().StringVarP(&metricDeleteName, "metric", "m", "", "metric name, alternative to id flag")
	metricDeleteCmd.MarkFlagRequired("metric")

	metricCmd.AddCommand(metricListCmd)
}
