package main

import (
	"fmt"
	"time"

	"github.com/demolemo/meth"
	"github.com/demolemo/meth/sqlite"

	"github.com/spf13/cobra"
)

var valueCmd = &cobra.Command{
	Use:   "value",
	Short: "Manage metric values",
}

var (
	valueAddID     int
	valueAddMetric string
	valueAddValue  float64
	valueAddAt     string
)

var valueAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a value to a metric",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		createdAt := time.Now()
		if cmd.Flags().Changed("at") {
			var err error
			createdAt, err = sqlite.ParseTime(valueAddAt)
			if err != nil {
				return err
			}
		}
		metricID, err := resolveMetricID(valueAddID, valueAddMetric)
		if err != nil {
			return err
		}
		mv := &meth.MetricValue{MetricID: metricID, Value: valueAddValue, CreatedAt: createdAt}
		// NOTE: change this in the next pass, this would not work later
		if _, err := sqlite.CreateMetricValue(db, mv); err != nil {
			return err
		}
		fmt.Printf("Metric value added nigga - metricID: %d, value: %.2f, createdAt: %s\n", mv.MetricID, mv.Value, mv.CreatedAt.Format(sqlite.TimeLayout))
		return nil
	},
}

// NOTE: later we can invent an interface that updates this in bulk
// much like in wtfgo

var (
	valueUpdateMetricID int
	valueUpdateMetric   string
	valueUpdateValueID  int
	valueUpdateValue    float64
	valueUpdateAt       string
)

var valueUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a value inside of a metric",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		metricID, err := resolveMetricID(valueUpdateMetricID, valueUpdateMetric)
		if err != nil {
			return err
		}
		updatedAt := time.Now()
		mv := &meth.MetricValue{ID: valueUpdateValueID, MetricID: metricID, Value: valueUpdateValue, UpdatedAt: updatedAt}
		if _, err := sqlite.UpdateMetricValue(ms, mv); err != nil {
			return err
		}
		fmt.Printf("Metric value update nigga - metricID: %d, value: %.2f, updatedAt: %s\n", mv.MetricID, mv.Value, mv.UpdatedAt.Format(sqlite.TimeLayout))
		return nil
	},
}

var (
	valueListMetricID int
	valueListMetric   string
	valueListLimit    int
)

var valueListCmd = &cobra.Command{
	Use:   "list",
	Short: "List values of a metric, newest first",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		metricID, err := resolveMetricID(valueListMetricID, valueListMetric)
		if err != nil {
			return err
		}
		mvs, err := sqlite.GetMetricValues(db, metricID, valueListLimit)
		if err != nil {
			return err
		}
		if len(mvs) == 0 {
			fmt.Printf("No values for metric %d\n", metricID)
			return nil
		}
		fmt.Printf("Values of metric %d:\n", metricID)
		for _, mv := range mvs {
			fmt.Printf("\tValue, ID: %d, Value: %g, CreatedAt: %s\n", mv.ID, mv.Value, mv.CreatedAt.Local().Format(sqlite.TimeLayout))
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
		rowsAffected, err := sqlite.DeleteMetricValueByID(db, valueDeleteValueID)
		if err != nil {
			return err
		}
		fmt.Printf("Value deleted, rows affected: %d\n", rowsAffected)
		return nil
	},
}

// --metric is an alternative to the metric id flag, cobra makes sure only one is set
func resolveMetricID(id int, name string) (int, error) {
	if name == "" {
		return id, nil
	}
	m, err := ms.GetMetricByName(name)
	if err != nil {
		return 0, err
	}
	if m == nil {
		return 0, fmt.Errorf("metric %q not found", name)
	}
	return m.ID, nil
}

//	var valueImportCmd = &cobra.Command{
//		Use:   "import",
//		Short: "Bulk import metric values",
//		Args:  cobra.NoArgs,
//		RunE: func(cmd *cobra.Command, args []string) error {
//			// first parse - this is a separete interface inside of csv package. return metricValues
//			// interface: accept path to the file, return metric values
//			// secondly - get an array of metrics and pass them to the writer
//			// interface writeMetrics(db, mvs) (int, err) where int is the number of errors that was written
//			rowsAffected, err := reader.readMetricValues(db)
//			rowsAffected, err := sqlite.DeleteMetricValueByID(db, valueDeleteValueID)
//			if err != nil {
//				return err
//			}
//			fmt.Printf("Value deleted, rows affected: %d\n", rowsAffected)
//			return nil
//		},
//	}
func init() {
	rootCmd.AddCommand(valueCmd)

	valueCmd.AddCommand(valueAddCmd)
	valueAddCmd.Flags().IntVar(&valueAddID, "id", 0, "metric ID")
	valueAddCmd.Flags().StringVarP(&valueAddMetric, "metric", "m", "", "metric name, alternative to --id")
	valueAddCmd.Flags().Float64VarP(&valueAddValue, "value", "v", 0, "value to record")
	valueAddCmd.Flags().StringVar(&valueAddAt, "at", "", `created_at, "2006-01-02 15:04" (default now)`)
	valueAddCmd.MarkFlagsOneRequired("id", "metric")
	valueAddCmd.MarkFlagsMutuallyExclusive("id", "metric")
	valueAddCmd.MarkFlagRequired("value")

	// NOTE: very very ugly interface, I know that yo.
	// i will make things better in the next pass
	valueCmd.AddCommand(valueUpdateCmd)
	valueUpdateCmd.Flags().IntVar(&valueUpdateMetricID, "mid", 0, "metric ID")
	valueUpdateCmd.Flags().StringVarP(&valueUpdateMetric, "metric", "m", "", "metric name, alternative to --mid")
	valueUpdateCmd.Flags().IntVar(&valueUpdateValueID, "vid", 0, "value ID")
	valueUpdateCmd.Flags().Float64VarP(&valueUpdateValue, "value", "v", 0, "value to record")
	valueUpdateCmd.MarkFlagsOneRequired("mid", "metric")
	valueUpdateCmd.MarkFlagsMutuallyExclusive("mid", "metric")
	valueUpdateCmd.MarkFlagRequired("vid")
	valueUpdateCmd.MarkFlagRequired("value")

	// NOTE: this interface could be changed to:
	// meth list - lists all of the metrics inside of the application
	// meth metric list - lists all the values (however, we have to pass the metric-id the either way)
	valueCmd.AddCommand(valueListCmd)
	valueListCmd.Flags().IntVar(&valueListMetricID, "id", 0, "metric ID")
	valueListCmd.Flags().StringVarP(&valueListMetric, "metric", "m", "", "metric name, alternative to --id")
	valueListCmd.Flags().IntVar(&valueListLimit, "limit", 5, "max values to show, -1 for all")
	valueListCmd.MarkFlagsOneRequired("id", "metric")
	valueListCmd.MarkFlagsMutuallyExclusive("id", "metric")

	valueCmd.AddCommand(valueDeleteCmd)
	valueDeleteCmd.Flags().IntVar(&valueDeleteValueID, "vid", 0, "value ID")
	valueDeleteCmd.MarkFlagRequired("vid")

	// valueCmd.AddCommand(valueImportCmd)
	// valueImportCmd.Flags().StringVarP(%)
}
