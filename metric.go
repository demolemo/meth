package meth

import (
	"errors"
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

type Metric struct {
	// Unique field that helps to identify this metric
	ID int `json:"id"`

	// Human readable name that helps to identify the metric further
	Name string `json:"name"`

	// Aggregation fields that define the aggregation rules inside of this metric
	AggRule     AggRule     `json:"aggRule"`
	AggInterval AggInterval `json:"aggInterval"`

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

type MetricService interface {
	// search metric by it's unique numeric ID, return ERRNOTFOUND
	// if there is no metric with according numeric id
	FindMetricByID(id int) (*Metric, error)

	// search metric by it's unique name, return ERRNOTFOUND
	// if there is no metric with according name yo
	FindMetricByName(name string) (*Metric, error)

	// pass a built metric so it could be recorded somewhere
	// bad comment, lack of an understanding
	CreateMetric(metric *Metric)

	// updates a given metric with new params, should we return new metric here? prolly not
	UpdateMetric(id int, upd MetricUpdate) (*Metric, error)

	// deletes a given metric and all values attached to it?
	// really good question to think about, are values deleted or retained?
	DeleteMetric(id int) error
}

type MetricUpdate struct {
	Name        string
	AggRule     AggRule
	AggInterval AggInterval
}
