package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/demolemo/meth"
	"github.com/demolemo/meth/parse"

	"github.com/spf13/cobra"
)

var valueCmd = &cobra.Command{
	Use:   "value",
	Short: "Manage metric values",
}

var (
	valueAddMetric string
	valueAddValue  float64
	valueAddAt     string
)

var valueAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a value to a metric",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		quietUsage(cmd)
		createdAt := time.Now()
		if cmd.Flags().Changed("at") {
			var err error
			createdAt, err = parse.ParseTime(valueAddAt)
			if err != nil {
				return err
			}
		}
		m, err := ms.GetMetric(valueAddMetric)
		if err != nil {
			return err
		}
		// FUCK, now we have to change the interface here
		mv := &meth.MetricValue{MetricID: m.ID, Value: valueAddValue, CreatedAt: createdAt}
		if _, err := vs.CreateMetricValue(mv); err != nil {
			return err
		}
		fmt.Printf("Metric value added nigga - metricID: %d, value: %.2f, createdAt: %s\n", mv.MetricID, mv.Value, mv.CreatedAt.Format(parse.TimeLayout))
		return nil
	},
}

// NOTE: later we can invent an interface that updates this in bulk
// much like in wtfgo

var (
	valueUpdateValueID int
	valueUpdateValue   float64
)

// the value id is enough to find the row, so this takes no metric. it used to
// resolve one out of two variables that were never bound to any flag, so every
// run failed with `metric "": not found`.
var valueUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a value by its ID",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		quietUsage(cmd)
		updatedAt := time.Now()
		mv := &meth.MetricValue{ID: valueUpdateValueID, Value: valueUpdateValue, UpdatedAt: updatedAt}
		if _, err := vs.UpdateMetricValue(mv); err != nil {
			return err
		}
		fmt.Printf("Metric value update nigga - valueID: %d, value: %.2f, updatedAt: %s\n", mv.ID, mv.Value, mv.UpdatedAt.Format(parse.TimeLayout))
		return nil
	},
}

var (
	valueListMetricName string
	valueListLimit      int
)

var valueListCmd = &cobra.Command{
	Use:   "list",
	Short: "List values of a metric, newest first",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		quietUsage(cmd)
		m, err := ms.GetMetric(valueListMetricName)
		if err != nil {
			return err
		}
		mvs, err := vs.GetMetricValues(m.ID, valueListLimit)
		if err != nil {
			return err
		}
		if len(mvs) == 0 {
			fmt.Printf("No values for metric %s\n", m.Name)
			return nil
		}
		fmt.Printf("Values of metric %s:\n", m.Name)
		for _, mv := range mvs {
			fmt.Printf("\tValue, ID: %d, Value: %g, CreatedAt: %s\n", mv.ID, mv.Value, mv.CreatedAt.Local().Format(parse.TimeLayout))
		}
		return nil
	},
}

var valueDeleteValueID int

var valueDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a value by its ID",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		quietUsage(cmd)
		rowsAffected, err := vs.DeleteMetricValue(valueDeleteValueID)
		if err != nil {
			return err
		}
		fmt.Printf("Value deleted, rows affected: %d\n", rowsAffected)
		return nil
	},
}

var (
	valueImportMetric string
	valueImportFile   string
)

var valueImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Bulk import metric values from a CSV file",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		quietUsage(cmd)
		m, err := ms.GetMetric(valueImportMetric)
		if err != nil {
			return err
		}
		// parse it all before touching the db - a bad last row imports nothing
		var mvs []*meth.MetricValue
		switch ext := strings.ToLower(filepath.Ext(valueImportFile)); ext {
		case ".csv":
			mvs, err = parse.CSV(valueImportFile, m.ID)
		case ".json":
			mvs, err = parse.HealthExportWeight(valueImportFile, m.ID)
		default:
			return fmt.Errorf("don't know how to read %q, want .csv or .json", ext)
		}
		if err != nil {
			return err
		}
		if len(mvs) == 0 {
			fmt.Printf("No values in %s, nothing imported\n", valueImportFile)
			return nil
		}
		res, err := vs.CreateMetricValues(mvs)
		if err != nil {
			return err
		}
		fmt.Printf("Imported %s into metric %s: %d new, %d updated, %d already there\n",
			valueImportFile, m.Name, res.Inserted, res.Updated, res.Unchanged)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(valueCmd)

	valueCmd.AddCommand(valueAddCmd)
	valueAddCmd.Flags().StringVarP(&valueAddMetric, "metric", "m", "", "metric name")
	valueAddCmd.Flags().Float64VarP(&valueAddValue, "value", "v", 0, "value to record")
	valueAddCmd.Flags().StringVar(&valueAddAt, "at", "", `created_at, "2006-01-02 15:04" (default now)`)
	valueAddCmd.MarkFlagRequired("metric")
	valueAddCmd.MarkFlagRequired("value")

	valueCmd.AddCommand(valueUpdateCmd)
	valueUpdateCmd.Flags().IntVar(&valueUpdateValueID, "id", 0, "value ID")
	valueUpdateCmd.Flags().Float64VarP(&valueUpdateValue, "value", "v", 0, "value to record")
	valueUpdateCmd.MarkFlagRequired("id")
	valueUpdateCmd.MarkFlagRequired("value")

	valueCmd.AddCommand(valueListCmd)
	valueListCmd.Flags().StringVarP(&valueListMetricName, "metric", "m", "", "metric name")
	valueListCmd.Flags().IntVar(&valueListLimit, "limit", 5, "max values to show, -1 for all")
	valueListCmd.MarkFlagRequired("metric")

	valueCmd.AddCommand(valueDeleteCmd)
	valueDeleteCmd.Flags().IntVar(&valueDeleteValueID, "id", 0, "value ID")
	valueDeleteCmd.MarkFlagRequired("id")

	valueCmd.AddCommand(valueImportCmd)
	valueImportCmd.Flags().StringVarP(&valueImportMetric, "metric", "m", "", "metric name")
	valueImportCmd.Flags().StringVarP(&valueImportFile, "file", "f", "", "file to import: .csv with value,created_at columns or a Health Export .json")
	valueImportCmd.MarkFlagRequired("metric")
	valueImportCmd.MarkFlagRequired("file")
}
