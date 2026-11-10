// Package meth is the domain: what a metric and a value are, and what any
// store of them has to be able to do. It depends on nothing but the standard
// library, so the storage and parsing packages can import it and it never
// imports them back.
package meth

type Metric struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MetricService is what sqlite.MetricService implements, asserted at compile
// time over there. The names say ByID and ByName because name is the handle a
// person uses and ids are what the database needs.
//
// Every method reports a miss the same way, by wrapping ErrNotFound: the
// lookups instead of returning a nil *Metric with a nil error, the writes
// instead of returning 0 rows affected and no error. So no caller has to
// nil-check or count rows to find out whether anything happened.
type MetricService interface {
	CreateMetric(m *Metric) error

	// case-insensitive, matching the unique index on metrics.name.
	// wraps ErrNotFound when no metric has that name
	GetMetric(name string) (*Metric, error)

	// returns <limit> number of metrics that exist in the service
	// metrics are sorted by an internal ID
	// NOTE: this could be an issue because after deletion
	// new metric could receive lower ID than the last stored metric
	// (sqlite works in that manner)
	GetMetrics(limit int) ([]*Metric, error)

	// wrap ErrNotFound when there was nothing to delete
	DeleteMetric(name string) error
}
