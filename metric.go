package metric_tracker

import (
	"fmt"
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
type Metric struct {
	ID int `json:"id"`
}
