package store

import (
	"database/sql"
	"errors"
	"time"
)

type Metric struct {
	ID   int
	Name string
}

type MetricUpdate struct {
	ID   int
	Name string
}

func CreateMetric(db *sql.DB, m *Metric) (int64, error) {
	sqlQuery := `INSERT INTO metrics VALUES($1, $2)`
	res, err := db.Exec(sqlQuery, m.ID, m.Name)
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
func GetMetricByID(db *sql.DB, id int) (*Metric, error) {
	sqlQuery := `
		SELECT * FROM metrics WHERE id = $id
	`
	rows, err := db.Query(sqlQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var name string

	for rows.Next() {
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		} else {
			return &Metric{ID: id, Name: name}, nil
		}
	}
	return nil, errors.New("err not found") // reuse this error in several places
}

func GetMetrics(db *sql.DB) ([]*Metric, error) {
	sqlQuery := `
		SELECT id, name FROM metrics
	`
	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mets []*Metric
	for rows.Next() {
		m := &Metric{}
		if err := rows.Scan(&m.ID, &m.Name); err != nil {
			return mets, err
		}
		mets = append(mets, m)
	}
	return mets, nil
}

func UpdateMetricByID(db *sql.DB, mu *MetricUpdate) (int64, error) {
	// do a check that metric exists first, for now we assume that it does
	sqlQuery := `UPDATE metrics SET name = $1 WHERE id = $2`
	res, err := db.Exec(sqlQuery, mu.Name, mu.ID)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

// solve the question with MetricValues. Do we need to delete values before the metric?
// are we working with values separately or we deleting them beforehand
func DeleteMetricByID(db *sql.DB, id int) (int64, error) {
	sqlQuery := `DELETE FROM metrics WHERE id = $id`
	res, err := db.Exec(sqlQuery, id)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

type MetricValue struct {
	ID        int
	MetricID  int
	Value     float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func CreateMetricValue(db *sql.DB, mv *MetricValue) (int64, error) {
	// we are worrying about corretly supplying the date outside of this function
	// this solution gives us a little bit more responsibility
	sqlQuery := `INSERT INTO metric_values(metric_id, value, created_at) VALUES($1, $2, $3)`
	res, err := db.Exec(sqlQuery, mv.MetricID, mv.Value, mv.CreatedAt)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func UpdateMetricValue(db *sql.DB, mv *MetricValue) (int64, error) {
	_, err := GetMetricByID(db, mv.MetricID)
	if err != nil {
		return 0, err
	}
	// NOTE: we can avoid checking here because the query will do implicit checking
	sqlQuery := `UPDATE metric_values SET value = $1, updated_at = $2 where id = $3 AND metric_id = $4`
	res, err := db.Exec(sqlQuery, mv.Value, mv.UpdatedAt, mv.ID, mv.MetricID)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

// weight and steps don't need exact seconds
const TimeLayout = "2006-01-02 15:04"

func ParseTime(s string) (time.Time, error) {
	return time.ParseInLocation(TimeLayout, s, time.Local)
}
