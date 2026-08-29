package meth

import (
	"errors"
	"fmt"
	"time"
	"unicode/utf8"
)

// two separate values
// Metric - represents some kind of an aggregate over the metric values (storage)
// MetricValue - singular metric value (more or less single)

// what should metric contain
// numeric value of the metric
// when was that metric created
// when was that metric updated
// last values of the metric (previous values of the metric)
// name of the metric

type AggRule int

const (
	AggMax AggRule = iota
	AggMin
	AggAvg
)

type AggInterval int

const (
	MinuteInterval AggInterval = iota
	HourInterval
	DayInterval
)

// Main structure of the program, responsible for storing metric values
type Metric struct {
	// Unique field that helps to identify this metric
	ID int `json:"id"`

	// Human readable name that helps to identify the metric further
	Name string `json:"name"`

	// Values that are stored inside of this metric
	Values *[]MetricValue `json:"values,omitempty"`
	// here i see the issue because this ensures that we need to create a new metric value for each aggregate
	// i don't know how to solve this issue yet
	AggValues *[]MetricValue `json:"aggValues,omitempty"`

	// Time fields, CreatedAt belongs to the Metric itself and UpdatedAt belongs to the underlying values
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (m *Metric) Validate() error {
	if m.Name == "" {
		return errors.New("Your name cannot be blank, nigga")
	} else if utf8.RuneCountInString(m.Name) > 120 {
		return errors.New("Your name cannot be that long, nigga")
	}
	return nil
}

// anything that implements this interface can access the metric
type MetricService interface {
	// search metric by it's unique numeric ID, return ERRNOTFOUND
	// if there is no metric with according numeric id
	FindMetricByID(id int) (*Metric, error)

	// search metric by it's unique name, return ERRNOTFOUND
	// if there is no metric with according name yo
	// NOTE: i'm not sure that we need this thing because this creates a hassle with maintaining
	// a set of metric names, for now it's commented out
	// FindMetricByName(name string) (*Metric, error)

	// pass a built metric so it could be recorded somewhere
	// where does the metric id come from? from some outer service?
	// do we validate that the name is unique on this stage?
	CreateMetric(metric *Metric)

	// generate a report with aggregated values
	GenerateMetricReport(id int, aggRule AggRule, aggInterval AggInterval) (MetricValueReport, error)

	// updates a given metric with new params
	UpdateMetric(id int, upd MetricUpdate) (*Metric, error)

	// deletes a given metric and all values attached to it?
	// really good question to think about, are values deleted or retained?
	DeleteMetric(id int) error
}

type MetricUpdate struct {
	Name string
}

// creating the aggregation values of the metric according to passed rules
// this solves the issue with storing aggregated values - we don't store them at all!
type MetricValueReport struct {
	// ID - later we can think about some kind of id field that will help caching requests with the same params
	// Name of the metric that produced that report
	Name string `json:"name"`

	// Rules which were applied to the underlying metric to arrive to this values
	// Hmmm, maybe we don't even need to store aggrule and agginterval near the metric.
	// metric is just raw values, that is implicit. i love that line of thinking

	// Aggregation fields that define the aggregation rules inside of this metric
	AggRule     AggRule     `json:"aggRule"`
	AggInterval AggInterval `json:"aggInterval"`

	// aggregated values according to aggregation rules above
	// is storing time a good idea here? maybe we can store ints instead?
	// we will try doing it that way first
	Records []MetricValueRecord `json:"aggValues"`
}

// MetricValueRecord represents an average metric value at a given point in time
// for the MetricValueReport.
// this is done so that we are not storing maps in there
type MetricValueRecord struct {
	Value     float32   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

// GoString prints a more easily readable representation for debugging.
// The timestamp field is represented as an RFC 3339 string instead of a pointer.
func (r *MetricValueRecord) GoString() string {
	return fmt.Sprintf("&meth.MetricValueRecord{Value:%d, Timestamp:%q}", r.Value, r.Timestamp.Format(time.RFC3339))
}
