package main

import (
	"time"
)

type MetricValue struct {
	// integer id that uniqely identifies this value
	ID int `json:"ID"`

	// value of the metric itself. float32 for now, don't see the point to make it more dense
	Value float32 `json:"Value"`

	// fields that help to identify the parent metric this value belongs to
	MetricID int `json:"MetricID"`
	Metric   *Metric

	// time related metrics
	CreatedAt time.Time
	UpdatedAt time.Time
}
