package meth

import (
	"time"
)

type Metric struct {
	ID   int
	Name string
}

type MetricUpdate struct {
	ID   int
	Name string
}

type MetricValue struct {
	ID        int
	MetricID  int
	Value     float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MetricService interface {
	// NOTE: we should remove int64 as param for all those methods, only viable inside of the sqlite
	CreateMetric(*Metric) (int64, error)

	// should throw error if metric with such ID doesn't exist
	GetMetricByID(id int) (*Metric, error)
	// retrieve metric by name, should throw an error if metric with such name was not found
	GetMetricByName(name string) (*Metric, error)

	// returns all metrics that exist in the service?
	// maybe we should add a limit here because this unbound behaviour is unsound
	GetMetrics() ([]*Metric, error)

	// should thorw error if metric with such ID doesn't exist
	UpdateMetricByID(mu *MetricUpdate) (int64, error)
	// UpdateMetricByID(mu *MetricUpdate) (*Metric, error) - this is how it's supposed to be in the future

	// delete metric by ID, should throw an error if metric with such ID was not found
	DeleteMetricByID(id int) (int64, error)
	// delete metric by name, should thorw an error if metric with such name was not found
	DeleteMetricByName(name string) (int64, error)
}
