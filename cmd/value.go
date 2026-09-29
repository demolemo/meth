package cmd

import (
	"fmt"
	"time"

	"meth/store"

	"github.com/spf13/cobra"
)

var valueCmd = &cobra.Command{
	Use:   "value",
	Short: "Manage metric values",
}

var (
	valueAddID    int
	valueAddValue float64
	valueAddAt    string
)

var valueAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a value to a metric",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		createdAt := time.Now()
		if cmd.Flags().Changed("at") {
			var err error
			createdAt, err = store.ParseTime(valueAddAt)
			if err != nil {
				return err
			}
		}
		mv := &store.MetricValue{MetricID: valueAddID, Value: valueAddValue, CreatedAt: createdAt}
		if _, err := store.CreateMetricValue(db, mv); err != nil {
			return err
		}
		fmt.Printf("Metric value added nigga - metricID: %d, value: %.2f, createdAt: %s\n", mv.MetricID, mv.Value, mv.CreatedAt.Format(store.TimeLayout))
		return nil
	},
}

// NOTE: later we can invent an interface that updates this in bulk
// much like in wtfgo

var (
	valueUpdateMetricID int
	valueUpdateValueID  int
	valueUpdateValue    float64
	valueUpdateAt       string
)

var valueUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a value inside of a metric",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		updatedAt := time.Now()
		mv := &store.MetricValue{ID: valueUpdateValueID, MetricID: valueUpdateMetricID, Value: valueAddValue, UpdatedAt: updatedAt}
		if _, err := store.UpdateMetricValue(db, mv); err != nil {
			return err
		}
		fmt.Printf("Metric value update nigga - metricID: %d, value: %.2f, createdAt: %s\n", mv.MetricID, mv.Value, mv.CreatedAt.Format(store.TimeLayout))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(valueCmd)

	valueCmd.AddCommand(valueAddCmd)
	valueAddCmd.Flags().IntVar(&valueAddID, "id", 0, "metric ID")
	valueAddCmd.Flags().Float64VarP(&valueAddValue, "value", "v", 0, "value to record")
	valueAddCmd.Flags().StringVar(&valueAddAt, "at", "", `created_at, "2006-01-02 15:04" (default now)`)
	valueAddCmd.MarkFlagRequired("id")
	valueAddCmd.MarkFlagRequired("value")

	// NOTE: very very ugly interface, I know that yo.
	// i will make things better in the next pass
	valueCmd.AddCommand(valueUpdateCmd)
	valueUpdateCmd.Flags().IntVar(&valueUpdateMetricID, "mid", 0, "metric ID")
	valueUpdateCmd.Flags().IntVar(&valueUpdateValueID, "vid", 0, "value ID")
	valueUpdateCmd.Flags().Float64VarP(&valueUpdateValue, "value", "v", 0, "value to record")
	valueUpdateCmd.MarkFlagRequired("mid")
	valueUpdateCmd.MarkFlagRequired("vid")
	valueUpdateCmd.MarkFlagRequired("v")
}
