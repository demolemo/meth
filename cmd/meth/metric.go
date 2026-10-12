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
	metricCreateUnit string
)

var metricCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a metric",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if strings.TrimSpace(metricCreateName) == "" {
			return fmt.Errorf("metric name can't be empty")
		}
		rowsAffected, err := ms.CreateMetric(&meth.Metric{ID: metricCreateID, Name: metricCreateName, Unit: metricCreateUnit})
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
		fmt.Printf("Metric, ID: %d, Name: %s, Unit: %s\n", met.ID, met.Name, met.Unit)
		return nil
	},
}

var (
	metricUpdateID   int
	metricUpdateName string
	metricUpdateUnit string
)

var metricUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Rename a metric",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// UpdateMetricByID writes both columns, so the one that was not
		// passed has to come from the row rather than from the zero value
		met, err := ms.GetMetricByID(metricUpdateID)
		if err != nil {
			return err
		}
		mu := &meth.MetricUpdate{ID: met.ID, Name: met.Name, Unit: met.Unit}
		if cmd.Flags().Changed("name") {
			if strings.TrimSpace(metricUpdateName) == "" {
				return fmt.Errorf("metric name can't be empty")
			}
			mu.Name = metricUpdateName
		}
		if cmd.Flags().Changed("unit") {
			mu.Unit = metricUpdateUnit
		}
		rowsAffected, err := ms.UpdateMetricByID(mu)
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

var (
	metricListPrefix string
	metricListLimit  int
)

var metricListCmd = &cobra.Command{
	Use:   "list",
	Short: "List metrics, optionally one dotted family",
	Long: `List metrics, optionally one dotted family.

A family is not a row anywhere, it is the start of a name:

  meth metric list --prefix sleep.       # every sleep metric
  meth metric list --prefix sleep.r      # sleep.rem and sleep.resp_rate`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		mets, err := ms.GetMetricsByPrefix(metricListPrefix)
		if err != nil {
			return err
		}
		if len(mets) == 0 {
			if metricListPrefix == "" {
				fmt.Printf("No metrics yet\n")
			} else {
				fmt.Printf("No metrics under %q\n", metricListPrefix)
			}
			return nil
		}

		shown := mets
		if metricListLimit > 0 && len(shown) > metricListLimit {
			shown = shown[:metricListLimit]
			fmt.Printf("%d metrics, first %d:\n", len(mets), len(shown))
		} else {
			fmt.Printf("%d metrics:\n", len(mets))
		}
		for _, m := range shown {
			fmt.Printf("\t%4d  %-22s %s\n", m.ID, m.Name, m.Unit)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(metricCmd)

	metricCmd.AddCommand(metricCreateCmd)
	metricCreateCmd.Flags().IntVar(&metricCreateID, "id", 0, "metric ID")
	metricCreateCmd.Flags().StringVarP(&metricCreateName, "name", "n", "", "metric name, a dotted path such as sleep.rem")
	metricCreateCmd.Flags().StringVarP(&metricCreateUnit, "unit", "u", "", `what the numbers are in, e.g. "h", "%", "kg"`)
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
	metricUpdateCmd.Flags().StringVarP(&metricUpdateUnit, "unit", "u", "", `new unit, e.g. "h", "%", "kg"`)
	metricUpdateCmd.MarkFlagRequired("id")
	metricUpdateCmd.MarkFlagsOneRequired("name", "unit")

	metricCmd.AddCommand(metricDeleteCmd)
	metricDeleteCmd.Flags().IntVar(&metricDeleteID, "id", 0, "metric ID")
	metricDeleteCmd.Flags().StringVarP(&metricDeleteName, "metric", "m", "", "metric name, alternative to id flag")
	metricDeleteCmd.MarkFlagsOneRequired("id", "metric")
	metricDeleteCmd.MarkFlagsMutuallyExclusive("id", "metric")

	metricCmd.AddCommand(metricListCmd)
	metricListCmd.Flags().StringVarP(&metricListPrefix, "prefix", "p", "", "only metrics whose name starts with this")
	metricListCmd.Flags().IntVar(&metricListLimit, "limit", 0, "max metrics to show, 0 for all")
}
