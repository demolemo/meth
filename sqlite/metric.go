package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/demolemo/meth"
)

var _ meth.MetricService = (*MetricService)(nil)

type MetricService struct {
	db *sql.DB
}

func NewMetricService(db *sql.DB) *MetricService {
	return &MetricService{db: db}
}

var ErrNotOneAffected = errors.New("Not one row affected after operation")

// here we should also update the metric under the metric pointer with ID
// if pointer arrives without an ID which should be expected
func (ms *MetricService) CreateMetric(m *meth.Metric) error {
	sqlQuery := `INSERT INTO metrics VALUES($1, $2)`
	res, err := ms.db.Exec(sqlQuery, m.ID, m.Name)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return ErrNotOneAffected
	}
	return nil
}

// case-insensitive, same as the unique index on metrics.name.
// wraps meth.ErrNotFound when no metric has that name
func (ms *MetricService) GetMetric(name string) (*meth.Metric, error) {
	sqlQuery := `SELECT id, name FROM metrics WHERE name = $1 COLLATE NOCASE`
	m := &meth.Metric{}
	err := ms.db.QueryRow(sqlQuery, name).Scan(&m.ID, &m.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("metric %q: %w", name, meth.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

// return stored metrics under a limit
// if limit == -1 then we return all metrics, caller behold
// (sqlite treats a negative LIMIT as no limit)
func (ms *MetricService) GetMetrics(limit int) ([]*meth.Metric, error) {
	// TODO: add sorting
	sqlQuery := `SELECT id, name FROM metrics LIMIT $1`
	rows, err := ms.db.Query(sqlQuery, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mets []*meth.Metric
	for rows.Next() {
		m := &meth.Metric{}
		if err := rows.Scan(&m.ID, &m.Name); err != nil {
			return mets, err
		}
		mets = append(mets, m)
	}

	if err := rows.Err(); err != nil {
		return mets, err
	}
	return mets, nil
}

// case-insensitive, matching GetMetricByName and the unique index.
// wraps meth.ErrNotFound when no metric has that name
func (ms *MetricService) DeleteMetric(name string) error {
	sqlQuery := `DELETE FROM metrics WHERE name = $1 COLLATE NOCASE`
	res, err := ms.db.Exec(sqlQuery, name)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("metric %q: %w", name, meth.ErrNotFound)
	}
	if rowsAffected != 1 {
		return ErrNotOneAffected
	}
	return nil
}
