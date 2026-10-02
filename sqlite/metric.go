package sqlite

import (
	"database/sql"
	"log"
	"time"

	"github.com/demolemo/meth"
)

type MetricService struct {
	db *sql.DB
}

func (ms *MetricService) FindMetricByID(id int) (*meth.Metric, error) {
	sqlStatement := `
		SELECT
			name,
			created_at,
			updated_at,
		FROM metrics 
		WHERE id == $1`
	row, err := ms.db.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	// this is supposed to trigger only one time. if it triggers
	// more than one time it means that we are fucked
	var (
		name      string
		createdAt time.Time
		updatedAt time.Time
	)
	for row.Next() {
		if err := row.Scan(&id, &name, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		log.Printf("FindMetricByID: %d %s %s %s\n", id, name, createdAt, updatedAt)
	}

	// are we creating the new value of the metric
	met := &meth.Metric{ID: id, Name: name, CreatedAt: createdAt, UpdatedAt: updatedAt}
	return met, nil
}

func (ms *MetricService) FindMetricByName(name string) (*meth.Metric, error) {
	return nil, nil
}

func (ms *MetricService) CreateMetric(metric *meth.Metric) error {
	return nil
}

func (ms *MetricService) UpdateMetric(id int, upd *meth.MetricUpdate) (*meth.Metric, error) {
	return nil, nil
}

func (ms *MetricService) DeleteMetric(id int) error {
	return nil
}
