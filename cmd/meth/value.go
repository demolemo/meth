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

const (
	dateLayout = "2006-01-02"

	// a reading anchored at a sleep's end lands on a second, not a minute,
	// so parse.TimeLayout would hide which of two sessions a row belongs to
	stampLayout = "2006-01-02 15:04:05 -07:00"
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
			createdAt, err = parse.ParseTime(valueAddAt)
			if err != nil {
				return err
			}
		}
		metricID, err := ms.ResolveMetricID(valueAddID, valueAddMetric)
		if err != nil {
			return err
		}
		mv := &meth.MetricValue{MetricID: metricID, Value: valueAddValue, CreatedAt: createdAt}
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

var valueUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a value inside of a metric",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// the metric comes from the value, not from a flag: UpdateMetricValue
		// matches on both and --id names a value
		old, err := vs.GetMetricValue(valueUpdateValueID)
		if err != nil {
			return err
		}
		updatedAt := time.Now()
		mv := &meth.MetricValue{ID: old.ID, MetricID: old.MetricID, Value: valueUpdateValue, UpdatedAt: updatedAt}
		if _, err := vs.UpdateMetricValue(mv); err != nil {
			return err
		}
		fmt.Printf("Metric value update nigga - metricID: %d, value: %.2f, updatedAt: %s\n", mv.MetricID, mv.Value, mv.UpdatedAt.Format(parse.TimeLayout))
		return nil
	},
}

var (
	valueListMetricID   int
	valueListMetricName string
	valueListLimit      int
)

var valueListCmd = &cobra.Command{
	Use:   "list",
	Short: "List values of a metric, newest first",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		metricID, err := ms.ResolveMetricID(valueListMetricID, valueListMetricName)
		if err != nil {
			return err
		}
		mvs, err := vs.GetMetricValues(metricID, valueListLimit)
		if err != nil {
			return err
		}
		met, err := ms.GetMetricByID(metricID)
		if err != nil {
			return err
		}
		if len(mvs) == 0 {
			fmt.Printf("No values for %s\n", met.Name)
			return nil
		}
		fmt.Printf("%s (%s), newest first:\n", met.Name, met.Unit)
		for _, mv := range mvs {
			fmt.Printf("\t%4d  %10.3f  %s\n", mv.ID, mv.Value, mv.CreatedAt.Local().Format(stampLayout))
		}
		return nil
	},
}

var (
	valueDayDate   string
	valueDayTZ     string
	valueDayPrefix string
)

var valueDayCmd = &cobra.Command{
	Use:   "day",
	Short: "Everything recorded on one calendar day, across every metric",
	Long: `Everything recorded on one calendar day, across every metric.

The day is a window over created_at, resolved in --tz. That zone matters:
a reading is one instant and the model has nowhere to keep the zone it was
recorded in, so a night slept in +08:00 and read back in +03:00 can fall on
the day before.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		loc, err := location(valueDayTZ)
		if err != nil {
			return err
		}
		from, err := time.ParseInLocation(dateLayout, valueDayDate, loc)
		if err != nil {
			return err
		}

		mvs, err := vs.GetValuesBetween(from, from.AddDate(0, 0, 1))
		if err != nil {
			return err
		}
		mets, err := ms.GetMetricsByPrefix(valueDayPrefix)
		if err != nil {
			return err
		}
		// one lookup table rather than a metric read per value: the join the
		// single values table does not do for us
		byID := make(map[int]*meth.Metric, len(mets))
		for _, m := range mets {
			byID[m.ID] = m
		}

		shown := 0
		for _, mv := range mvs {
			m, ok := byID[mv.MetricID]
			if !ok {
				continue
			}
			if shown == 0 {
				fmt.Printf("%s in %s:\n", valueDayDate, loc)
			}
			shown++
			fmt.Printf("\t%-22s %10.3f %-5s %s\n", m.Name, mv.Value, m.Unit, mv.CreatedAt.In(loc).Format(stampLayout))
		}
		if shown == 0 {
			fmt.Printf("Nothing recorded on %s in %s\n", valueDayDate, loc)
		}
		return nil
	},
}

// location reads an IANA name or a bare ±HH:MM offset, the form a whoop dump
// carries. An empty name is the machine's zone.
func location(name string) (*time.Location, error) {
	if name == "" {
		return time.Local, nil
	}
	if loc, err := time.LoadLocation(name); err == nil {
		return loc, nil
	}
	t, err := time.Parse("-07:00", name)
	if err != nil {
		return nil, fmt.Errorf("timezone %q: not an IANA name and not a ±HH:MM offset", name)
	}
	_, offset := t.Zone()
	return time.FixedZone(name, offset), nil
}

var valueDeleteValueID int

var valueDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a value by its ID",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		rowsAffected, err := vs.DeleteMetricValue(valueDeleteValueID)
		if err != nil {
			return err
		}
		fmt.Printf("Value deleted, rows affected: %d\n", rowsAffected)
		return nil
	},
}

var (
	valueImportID     int
	valueImportMetric string
	valueImportFile   string
)

var valueImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Bulk import metric values from a CSV file",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		metricID, err := ms.ResolveMetricID(valueImportID, valueImportMetric)
		if err != nil {
			return err
		}
		// parse it all before touching the db - a bad last row imports nothing
		var mvs []*meth.MetricValue
		switch ext := strings.ToLower(filepath.Ext(valueImportFile)); ext {
		case ".csv":
			mvs, err = parse.CSV(valueImportFile, metricID)
		case ".json":
			mvs, err = parse.HealthExportWeight(valueImportFile, metricID)
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
		fmt.Printf("Imported %s into metric %d: %d new, %d updated, %d already there\n",
			valueImportFile, metricID, res.Inserted, res.Updated, res.Unchanged)
		return nil
	},
}

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
	valueUpdateCmd.Flags().IntVar(&valueUpdateValueID, "id", 0, "value ID")
	valueUpdateCmd.Flags().Float64VarP(&valueUpdateValue, "value", "v", 0, "value to record")
	valueUpdateCmd.MarkFlagRequired("id")
	valueUpdateCmd.MarkFlagRequired("value")

	valueCmd.AddCommand(valueListCmd)
	valueListCmd.Flags().IntVar(&valueListMetricID, "id", 0, "metric ID")
	valueListCmd.Flags().StringVarP(&valueListMetricName, "metric", "m", "", "metric name, alternative to --id")
	valueListCmd.Flags().IntVar(&valueListLimit, "limit", 5, "max values to show, -1 for all")
	valueListCmd.MarkFlagsOneRequired("id", "metric")
	valueListCmd.MarkFlagsMutuallyExclusive("id", "metric")

	valueCmd.AddCommand(valueDayCmd)
	valueDayCmd.Flags().StringVarP(&valueDayDate, "date", "d", "", "the day to show, 2006-01-02")
	valueDayCmd.Flags().StringVar(&valueDayTZ, "tz", "", "zone the day is resolved in: an IANA name or ±HH:MM (default local)")
	valueDayCmd.Flags().StringVarP(&valueDayPrefix, "prefix", "p", "", "only metrics whose name starts with this")
	valueDayCmd.MarkFlagRequired("date")

	valueCmd.AddCommand(valueDeleteCmd)
	valueDeleteCmd.Flags().IntVar(&valueDeleteValueID, "id", 0, "value ID")
	valueDeleteCmd.MarkFlagRequired("id")

	valueCmd.AddCommand(valueImportCmd)
	valueImportCmd.Flags().IntVar(&valueImportID, "id", 0, "metric ID")
	valueImportCmd.Flags().StringVarP(&valueImportMetric, "metric", "m", "", "metric name, alternative to --id")
	valueImportCmd.Flags().StringVarP(&valueImportFile, "file", "f", "", "file to import: .csv with value,created_at columns or a Health Export .json")
	valueImportCmd.MarkFlagsOneRequired("id", "metric")
	valueImportCmd.MarkFlagsMutuallyExclusive("id", "metric")
	valueImportCmd.MarkFlagRequired("file")
}
