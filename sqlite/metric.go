package sqlite

import (
	"database/sql"
)

type MetricService struct {
	db *sql.DB
}

func (ms *MetricService) FindMetricByID(id int) (*Metric, error) {
	return nil, nil
}

func (ms *MetricService) FindMetricByName(name string) (*Metric, error) {
	return nil, nil
}

func (ms *MetricService) CreateMetric(metric *Metric) error {
	return nil
}

func (ms *MetricService) UpdateMetric(id int, upd MetricUpdate) (*Metric, error) {
	return nil, nil
}

func (ms *MetricService) DeleteMetric(id int) error {
	return nil
}
