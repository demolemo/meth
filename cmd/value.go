package cmd

import (
	"fmt"
	"time"

	"checking/store"

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

func init() {
	rootCmd.AddCommand(valueCmd)

	valueCmd.AddCommand(valueAddCmd)
	valueAddCmd.Flags().IntVar(&valueAddID, "id", 0, "metric ID")
	valueAddCmd.Flags().Float64VarP(&valueAddValue, "value", "v", 0, "value to record")
	valueAddCmd.Flags().StringVar(&valueAddAt, "at", "", `created_at, "2006-01-02 15:04" (default now)`)
	valueAddCmd.MarkFlagRequired("id")
	valueAddCmd.MarkFlagRequired("value")
}
