package meth

type Metric struct {
	ID   int
	Name string
}

type MetricUpdate struct {
	ID   int
	Name string
}

type MetricService interface {
	// NOTE: we should remove int64 as param for all those methods, only viable inside of the sqlite
	CreateMetric(*Metric) (int64, error)

	// should throw error if metric with such ID doesn't exist
	GetMetric(id int) (*Metric, error)

	// returns all metrics that exist in the service?
	// maybe we should add a limit here because this unbound behaviour is unsound
	GetMetrics() ([]*Metric, error)

	// should thorw error if metric with such ID doesn't exist
	UpdateMetric(mu *MetricUpdate) (int64, error)
	// UpdateMetricByID(mu *MetricUpdate) (*Metric, error) - this is how it's supposed to be in the future

	// delete metric by ID, should throw an error if metric with such ID was not found
	DeleteMetric(id int) (int64, error)

	// resolve metric id to not duplicate all of the functionality by name
	ResolveMetricID(id int, name string) (int, error)
}
