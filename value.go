package meth

import "time"

type MetricValue struct {
	ID        int
	MetricID  int
	Value     float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ImportResult counts what a batch did. Lives here rather than in sqlite
// because ValueService names it and sqlite already imports this package.
type ImportResult struct {
	Inserted  int
	Updated   int
	Unchanged int
}

// ValueService is what sqlite.ValueService implements, asserted at compile
// time over there.
//
// NOTE: remove all ints from the interface
type ValueService interface {
	// CRUD assemble
	// THIS IS CCCC
	CreateMetricValue(mv *MetricValue) (int64, error)

	// a batch in one transaction, for imports
	CreateMetricValues(mvs []*MetricValue) (ImportResult, error)

	// THIS IS R BUT IT'S A G because im retarded
	// wraps ErrNotFound when no value has that id
	GetMetricValue(valueID int) (*MetricValue, error)
	// THIS IS R BUT IT'S a G again
	GetMetricValues(metricID int, limit int) ([]*MetricValue, error)
	// this is U no surprise this time
	// keyed on mv.ID alone, mv.MetricID is ignored.
	// wraps ErrNotFound when no value has that id
	UpdateMetricValue(mv *MetricValue) (int64, error)
	// this is D. don't get any funny ideas here
	// wraps ErrNotFound when there was nothing to delete
	DeleteMetricValue(valueID int) (int64, error)
}
