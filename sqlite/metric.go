package sqlite

import (
	"database/sql"
	"errors"
	"time"

	"github.com/demolemo/meth"
)

type MetricService struct {
	db *sql.DB
}

func NewMetricService(db *sql.DB) *MetricService {
	return &MetricService{db: db}
}

func (ms *MetricService) CreateMetric(m *meth.Metric) (int64, error) {
	sqlQuery := `INSERT INTO metrics VALUES($1, $2)`
	res, err := ms.db.Exec(sqlQuery, m.ID, m.Name)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

// returns nil, nil when there is no metric with that id
func (ms *MetricService) GetMetricByID(id int) (*meth.Metric, error) {
	sqlQuery := `
		SELECT * FROM metrics WHERE id = $id
	`
	rows, err := ms.db.Query(sqlQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var name string

	for rows.Next() {
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		} else {
			return &meth.Metric{ID: id, Name: name}, nil
		}
	}
	return nil, errors.New("err not found") // reuse this error in several places
}

// case-insensitive, same as the unique index on metrics.name
// returns nil, nil when there is no metric with that name
func (ms *MetricService) GetMetricByName(name string) (*meth.Metric, error) {
	sqlQuery := `SELECT id, name FROM metrics WHERE name = $1 COLLATE NOCASE`
	m := &meth.Metric{}
	err := ms.db.QueryRow(sqlQuery, name).Scan(&m.ID, &m.Name)
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

// should the ID be primary key that we can affect? the name should be identifier too
// return *Metric here
func (ms *MetricService) UpdateMetricByID(mu *meth.MetricUpdate) (int64, error) {
	// do a check that metric exists first, for now we assume that it does
	sqlQuery := `UPDATE metrics SET name = $1 WHERE id = $2`
	res, err := ms.db.Exec(sqlQuery, mu.Name, mu.ID)
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
	sqlQuery := `DELETE FROM metrics WHERE id = $id`
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
	sqlQuery := `DELETE FROM metrics WHERE name = $name`
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
	return m.ID, nil
}

// weight and steps don't need exact seconds
// NOTE: this bullshit needs to be exported better or not exported at all.
// additional dep, bad stuff
const TimeLayout = "2006-01-02 15:04"

func ParseTime(s string) (time.Time, error) {
	return time.ParseInLocation(TimeLayout, s, time.Local)
}
