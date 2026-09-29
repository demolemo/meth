package cmd

import (
	"fmt"
	"strconv"
	"time"

	"checking/store"

	"github.com/spf13/cobra"
)

var valueAddCmd = &cobra.Command{
	Use:   `add ID value ["2006-01-02 15:04"]`,
	Short: "Add a value to a metric (created_at defaults to now)",
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		value, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return err
		}
		createdAt := time.Now()
		if len(args) == 3 {
			createdAt, err = store.ParseTime(args[2])
			if err != nil {
				return err
			}
		}
		mv := &store.MetricValue{MetricID: id, Value: value, CreatedAt: createdAt}
		if _, err := store.CreateMetricValue(db, mv); err != nil {
			return err
		}
		fmt.Printf("Metric value added nigga - metricID: %d, value: %.2f, createdAt: %s\n", mv.MetricID, mv.Value, mv.CreatedAt.Format(store.TimeLayout))
		return nil
	},
}

func init() {
	valueCmd.AddCommand(valueAddCmd)
}
