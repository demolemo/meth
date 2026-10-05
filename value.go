package meth

import "time"

type MetricValue struct {
	ID        int
	MetricID  int
	Value     float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NOTE: remove all ints from the interface
type ValueService interface {
	// CRUD assemble
	// THIS IS CCCC
	CreateMetricValue(mv *MetricValue) (int64, error)
	// THIS IS R BUT IT'S A G because im retarded
	GetMetricValue(valueID int) ([]*MetricValue, error)
	// THIS IS R BUT IT'S a G again
	GetMetricValues(metricID int, limit int) ([]*MetricValue, error)
	// this is U no surprise this time
	UpdateMetricValue(mv *MetricValue) (int64, error)
	// this is D. don't get any funny ideas here
	DeleteMetricValue(valueID int) (int64, error)
}
