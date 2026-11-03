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
	if rowsAffected != 0 {
		return ErrNotOneAffected
	}
	return nil
}

// wraps meth.ErrNotFound when no metric has that id
func (ms *MetricService) GetMetricByID(id int) (*meth.Metric, error) {
	sqlQuery := `SELECT id, name FROM metrics WHERE id = $1`
	m := &meth.Metric{}
	err := ms.db.QueryRow(sqlQuery, id).Scan(&m.ID, &m.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("metric %d: %w", id, meth.ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

// case-insensitive, same as the unique index on metrics.name.
// wraps meth.ErrNotFound when no metric has that name
func (ms *MetricService) GetMetricByName(name string) (*meth.Metric, error) {
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

// return all stored metric. should i add limit here?
func (ms *MetricService) GetMetrics() ([]*meth.Metric, error) {
	sqlQuery := `
		SELECT id, name FROM metrics
	`
	rows, err := ms.db.Query(sqlQuery)
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
	return mets, nil
}

// wraps meth.ErrNotFound when no metric has that id. sqlite counts a row as
// changed even when the new name equals the old one, so renaming a metric to
// what it is already still reports 1 and not a spurious miss.
//
// should the ID be primary key that we can affect? the name should be identifier too
// return *Metric here
func (ms *MetricService) UpdateMetricByID(mu *meth.MetricUpdate) error {
	sqlQuery := `UPDATE metrics SET name = $1 WHERE id = $2`
	res, err := ms.db.Exec(sqlQuery, mu.Name, mu.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("metric %d: %w", mu.ID, meth.ErrNotFound)
	}
	if rowsAffected != 1 {
		return ErrNotOneAffected
	}
	return nil
}

// wraps meth.ErrNotFound when no metric has that id.
//
// a metric that still has values fails on the foreign key instead, which is
// the open question in the README: do values get deleted first, or is this
// supposed to cascade?
func (ms *MetricService) DeleteMetricByID(id int) error {
	sqlQuery := `DELETE FROM metrics WHERE id = $1`
	res, err := ms.db.Exec(sqlQuery, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("metric %d: %w", id, meth.ErrNotFound)
	}
	// this logic might be a little bit redundant
	if rowsAffected != 1 {
		return ErrNotOneAffected
	}
	return nil
}

// case-insensitive, matching GetMetricByName and the unique index.
// wraps meth.ErrNotFound when no metric has that name
func (ms *MetricService) DeleteMetricByName(name string) error {
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

// a non-zero id wins, otherwise the name is looked up. Either way the metric
// has to exist: every caller is about to read or write its values, and an
// unchecked id surfaces later as "FOREIGN KEY constraint failed" or as a
// delete that quietly matches no rows.
func (ms *MetricService) ResolveMetricID(id int, name string) (int, error) {
	if id != 0 {
		if _, err := ms.GetMetricByID(id); err != nil {
			return 0, err
		}
		return id, nil
	}

	m, err := ms.GetMetricByName(name)
	if err != nil {
		return 0, err
	}
	return m.ID, nil
}
