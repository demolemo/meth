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

func (ms *MetricService) CreateMetric(m *meth.Metric) (int64, error) {
	sqlQuery := `INSERT INTO metrics(id, name, unit) VALUES($1, $2, $3)`
	res, err := ms.db.Exec(sqlQuery, m.ID, m.Name, m.Unit)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

// EnsureMetric returns the metric called name, creating it with unit when no
// metric has that name yet. An existing metric keeps the unit it already has,
// unless that unit contradicts the one asked for: filing brpm readings under
// a metric someone set to kg is worse than a failed import.
func (ms *MetricService) EnsureMetric(name, unit string) (*meth.Metric, error) {
	m, err := ms.GetMetricByName(name)
	if err != nil {
		return nil, err
	}
	if m != nil {
		if m.Unit != "" && unit != "" && m.Unit != unit {
			return nil, fmt.Errorf("metric %q is in %s, not %s", m.Name, m.Unit, unit)
		}
		return m, nil
	}

	// no id, so sqlite assigns one. CreateMetric takes the id from the caller
	// and an unset one would land a metric at row 0
	sqlQuery := `INSERT INTO metrics(name, unit) VALUES($1, $2)`
	res, err := ms.db.Exec(sqlQuery, name, unit)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &meth.Metric{ID: int(id), Name: name, Unit: unit}, nil
}

// returns an error when there is no metric with that id.
// NOTE: GetMetricByName returns nil, nil for the same situation. two
// different not-found contracts on one service, pick one and make it a
// sentinel error both can return
func (ms *MetricService) GetMetricByID(id int) (*meth.Metric, error) {
	sqlQuery := `SELECT id, name, unit FROM metrics WHERE id = $1`
	m := &meth.Metric{}
	err := ms.db.QueryRow(sqlQuery, id).Scan(&m.ID, &m.Name, &m.Unit)
	if err == sql.ErrNoRows {
		return nil, errors.New("err not found") // reuse this error in several places
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

// case-insensitive, same as the unique index on metrics.name
// returns nil, nil when there is no metric with that name
func (ms *MetricService) GetMetricByName(name string) (*meth.Metric, error) {
	sqlQuery := `SELECT id, name, unit FROM metrics WHERE name = $1 COLLATE NOCASE`
	m := &meth.Metric{}
	err := ms.db.QueryRow(sqlQuery, name).Scan(&m.ID, &m.Name, &m.Unit)
	if err == sql.ErrNoRows {
		// TODO: return an error not found here
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

// return all stored metric. should i add limit here?
func (ms *MetricService) GetMetrics() ([]*meth.Metric, error) {
	return ms.GetMetricsByPrefix("")
}

// GetMetricsByPrefix is how a dotted family is read back: `sleep.` matches
// every sleep metric because their names start with those six characters and
// for no other reason.
//
// substr rather than LIKE, so a `_` or a `%` in the prefix is six characters
// and not a wildcard.
func (ms *MetricService) GetMetricsByPrefix(prefix string) ([]*meth.Metric, error) {
	sqlQuery := `
		SELECT id, name, unit FROM metrics
		WHERE substr(name, 1, length($1)) = $1 COLLATE NOCASE
		ORDER BY name COLLATE NOCASE
	`
	rows, err := ms.db.Query(sqlQuery, prefix)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mets []*meth.Metric
	for rows.Next() {
		m := &meth.Metric{}
		if err := rows.Scan(&m.ID, &m.Name, &m.Unit); err != nil {
			return mets, err
		}
		mets = append(mets, m)
	}
	return mets, rows.Err()
}

// should the ID be primary key that we can affect? the name should be identifier too
// return *Metric here
func (ms *MetricService) UpdateMetricByID(mu *meth.MetricUpdate) (int64, error) {
	// do a check that metric exists first, for now we assume that it does
	sqlQuery := `UPDATE metrics SET name = $1, unit = $2 WHERE id = $3`
	res, err := ms.db.Exec(sqlQuery, mu.Name, mu.Unit, mu.ID)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

// solve the question with meth.MetricValues. Do we need to delete values before the metric?
// are we working with values separately or we deleting them beforehand
func (ms *MetricService) DeleteMetricByID(id int) (int64, error) {
	sqlQuery := `DELETE FROM metrics WHERE id = $1`
	res, err := ms.db.Exec(sqlQuery, id)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func (ms *MetricService) DeleteMetricByName(name string) (int64, error) {
	sqlQuery := `DELETE FROM metrics WHERE name = $1 COLLATE NOCASE`
	res, err := ms.db.Exec(sqlQuery, name)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func (ms *MetricService) ResolveMetricID(id int, name string) (int, error) {
	if id != 0 {
		return id, nil
	}

	m, err := ms.GetMetricByName(name)
	if err != nil {
		return 0, err
	}
	if m == nil {
		return 0, fmt.Errorf("metric %q not found", name)
	}
	return m.ID, nil
}
