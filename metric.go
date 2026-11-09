// Package meth is the domain: what a metric and a value are, and what any
// store of them has to be able to do. It depends on nothing but the standard
// library, so the storage and parsing packages can import it and it never
// imports them back.
package meth

type Metric struct {
	ID   int `json:"id"`
	Name string `json:"name"`
}

type MetricUpdate struct {
	ID   int
	Name string
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

	// wraps ErrNotFound when no metric has that id
	GetMetricByID(id int) (*Metric, error)

	// case-insensitive, matching the unique index on metrics.name.
	// wraps ErrNotFound when no metric has that name
	GetMetricByName(name string) (*Metric, error)

	// returns all metrics that exist in the service?
	// maybe we should add a limit here because this unbound behaviour is unsound
	GetMetrics() ([]*Metric, error)

	// wraps ErrNotFound when no metric has that id, so a rename that matched
	// nothing is an error and not a report of 0 rows affected
	UpdateMetricByID(mu *MetricUpdate) error

	// both wrap ErrNotFound when there was nothing to delete
	DeleteMetricByID(id int) error
	DeleteMetricByName(name string) error

	// resolve metric id to not duplicate all of the functionality by name.
	// a non-zero id wins, otherwise the name is looked up. wraps ErrNotFound
	// if the metric is not there either way
	ResolveMetricID(id int, name string) (int, error)
}
