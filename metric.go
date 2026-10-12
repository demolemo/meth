// Package meth is the domain: what a metric and a value are, and what any
// store of them has to be able to do. It depends on nothing but the standard
// library, so the storage and parsing packages can import it and it never
// imports them back.
package meth

type Metric struct {
	ID int

	// a dotted path: "sleep.rem". NOCASE unique, and the whole of the
	// nesting - there is no parent row, grouping a family is a prefix query.
	Name string

	// what the numbers are in: "h", "%", "count", "brpm", "bpm", "ms",
	// "degC", "kg". value is a bare float, so without this a kilogram and a
	// pound are the same column.
	Unit string
}

type MetricUpdate struct {
	ID   int
	Name string
	Unit string
}

// MetricService is what sqlite.MetricService implements, asserted at compile
// time over there. The names say ByID and ByName because name is the handle a
// person uses and ids are what the database needs.
//
// NOTE: we should remove int64 as param for all those methods, only viable inside of the sqlite
type MetricService interface {
	CreateMetric(m *Metric) (int64, error)

	// creates the metric when the name is free, otherwise returns the one
	// already there. an importer invents a couple of dozen names and should
	// not care which of them it is the first to see.
	EnsureMetric(name, unit string) (*Metric, error)

	// returns an error if no metric has that id
	GetMetricByID(id int) (*Metric, error)

	// case-insensitive, matching the unique index on metrics.name.
	// returns nil, nil when no metric has that name
	GetMetricByName(name string) (*Metric, error)

	// returns all metrics that exist in the service?
	// maybe we should add a limit here because this unbound behaviour is unsound
	GetMetrics() ([]*Metric, error)

	// every metric whose name starts with prefix, case-insensitive, ordered
	// by name. an empty prefix is every metric. this is the only way a family
	// is addressable: "sleep." is not a thing in the database, it is a string
	// that happens to start a lot of names.
	GetMetricsByPrefix(prefix string) ([]*Metric, error)

	UpdateMetricByID(mu *MetricUpdate) (int64, error)

	DeleteMetricByID(id int) (int64, error)
	DeleteMetricByName(name string) (int64, error)

	// resolve metric id to not duplicate all of the functionality by name.
	// a non-zero id wins, otherwise the name is looked up
	ResolveMetricID(id int, name string) (int, error)
}
