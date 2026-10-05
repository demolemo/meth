package sqlite

import (
	"database/sql"

	"github.com/demolemo/meth"
)

type ValueService struct {
	db *sql.DB
}

func NewValueService(db *sql.DB) *ValueService {
	return &ValueService{db: db}
}

func (vs *ValueService) CreateMetricValue(mv *meth.MetricValue) (int64, error) {
	// we are worrying about corretly supplying the date outside of this function
	// this solution gives us a little bit more responsibility
	sqlQuery := `INSERT INTO metric_values(metric_id, value, created_at) VALUES($1, $2, $3)`
	res, err := vs.db.Exec(sqlQuery, mv.MetricID, mv.Value, mv.CreatedAt)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func (vs *ValueService) GetMetricValue(valueID int) (*meth.MetricValue, error) {
	sqlQuery := `
		SELECT id, metric_id, value, created_at, updated_at
		FROM metric_values
		WHERE id = $id
	`
	rows, err := vs.db.Query(sqlQuery, valueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mv := &meth.MetricValue{}
	for rows.Next() {
		if err := rows.Scan(&mv.ID, &mv.MetricID, &mv.Value, &mv.CreatedAt, &mv.UpdatedAt); err != nil {
			return nil, err
		}
	}
	return mv, nil
}

// newest first. limit < 0 means no limit
// NOTE: add sorting rules
func (vs *ValueService) GetMetricValues(metricID int, limit int) ([]*meth.MetricValue, error) {
	sqlQuery := `
		SELECT id, metric_id, value, created_at, updated_at FROM metric_values
		WHERE metric_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	rows, err := vs.db.Query(sqlQuery, metricID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mvs []*meth.MetricValue
	for rows.Next() {
		mv := &meth.MetricValue{}
		if err := rows.Scan(&mv.ID, &mv.MetricID, &mv.Value, &mv.CreatedAt, &mv.UpdatedAt); err != nil {
			return mvs, err
		}
		mvs = append(mvs, mv)
	}
	return mvs, rows.Err()
}

// QUESTION: should we send a value update instead of the value itself?
func (vs *ValueService) UpdateMetricValue(mv *meth.MetricValue) (int64, error) {
	sqlQuery := `UPDATE metric_values SET value = $1, updated_at = $2 where id = $3 AND metric_id = $4`
	res, err := vs.db.Exec(sqlQuery, mv.Value, mv.UpdatedAt, mv.ID, mv.MetricID)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func (vs *ValueService) DeleteMetricValue(valueID int) (int64, error) {
	sqlQuery := `DELETE FROM metric_values WHERE id = $1`
	res, err := vs.db.Exec(sqlQuery, valueID)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}
