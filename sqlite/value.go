package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/demolemo/meth"
)

// the interface in the meth package is the contract, this line is what
// makes a rename here a build failure instead of silent drift
var _ meth.ValueService = (*ValueService)(nil)

type ValueService struct {
	db *sql.DB
}

func NewValueService(db *sql.DB) *ValueService {
	return &ValueService{db: db}
}

func (vs *ValueService) CreateMetricValue(mv *meth.MetricValue) (int64, error) {
	// we are worrying about corretly supplying the date outside of this function
	// this solution gives us a little bit more responsibility
	sqlQuery := `INSERT INTO metric_values(metric_id, value, created_at, updated_at) VALUES($1, $2, $3, $4)`
	res, err := vs.db.Exec(sqlQuery, mv.MetricID, mv.Value, mv.CreatedAt, time.Now())
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

// wraps meth.ErrNotFound when no value has that id. It used to hand back a
// zero-valued MetricValue and a nil error, which reads as "value 0 of metric
// 0 is 0" at every call site.
func (vs *ValueService) GetMetricValue(valueID int) (*meth.MetricValue, error) {
	sqlQuery := `
		SELECT id, metric_id, value, created_at, updated_at
		FROM metric_values
		WHERE id = $1
	`
	mv := &meth.MetricValue{}
	err := vs.db.QueryRow(sqlQuery, valueID).Scan(&mv.ID, &mv.MetricID, &mv.Value, &mv.CreatedAt, &mv.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("value %d: %w", valueID, meth.ErrNotFound)
	}
	if err != nil {
		return nil, err
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

// keyed on the value id alone. mv.MetricID is not in the WHERE clause: id is
// the primary key, so scoping by metric could never match a different row, and
// requiring it meant a caller had to name the metric just to edit one reading.
//
// wraps meth.ErrNotFound when no value has that id
//
// QUESTION: should we send a value update instead of the value itself?
func (vs *ValueService) UpdateMetricValue(mv *meth.MetricValue) (int64, error) {
	sqlQuery := `UPDATE metric_values SET value = $1, updated_at = $2 WHERE id = $3`
	res, err := vs.db.Exec(sqlQuery, mv.Value, mv.UpdatedAt, mv.ID)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	if rowsAffected == 0 {
		return 0, fmt.Errorf("value %d: %w", mv.ID, meth.ErrNotFound)
	}
	return rowsAffected, nil
}

// wraps meth.ErrNotFound when no value has that id, rather than reporting a
// successful delete of nothing
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
	if rowsAffected == 0 {
		return 0, fmt.Errorf("value %d: %w", valueID, meth.ErrNotFound)
	}
	return rowsAffected, nil
}

// CreateMetricValues writes a batch in one transaction, treating
// (metric_id, created_at) as a reading's identity: re-importing an overlapping
// export changes nothing unless the number did. Two statements rather than an
// upsert because no unique index backs that pair.
func (vs *ValueService) CreateMetricValues(mvs []*meth.MetricValue) (meth.ImportResult, error) {
	var res meth.ImportResult

	tx, err := vs.db.Begin()
	if err != nil {
		return res, err
	}
	defer tx.Rollback()

	update, err := tx.Prepare(`
		UPDATE metric_values SET value = $1, updated_at = $2
		WHERE metric_id = $3 AND created_at = $4 AND value <> $5
	`)
	if err != nil {
		return res, err
	}
	defer update.Close()

	insert, err := tx.Prepare(`
		INSERT INTO metric_values(metric_id, value, created_at, updated_at)
		SELECT $1, $2, $3, $4
		WHERE NOT EXISTS (
			SELECT 1 FROM metric_values WHERE metric_id = $5 AND created_at = $6
		)
	`)
	if err != nil {
		return res, err
	}
	defer insert.Close()

	now := time.Now()
	for i, mv := range mvs {
		changed, err := rowsAffected(update.Exec(mv.Value, now, mv.MetricID, mv.CreatedAt, mv.Value))
		if err != nil {
			return res, fmt.Errorf("value %d of %d: %w", i+1, len(mvs), err)
		}
		if changed > 0 {
			res.Updated++
			continue
		}

		added, err := rowsAffected(insert.Exec(mv.MetricID, mv.Value, mv.CreatedAt, now, mv.MetricID, mv.CreatedAt))
		if err != nil {
			return res, fmt.Errorf("value %d of %d: %w", i+1, len(mvs), err)
		}
		if added > 0 {
			res.Inserted++
		} else {
			res.Unchanged++
		}
	}
	return res, tx.Commit()
}

func rowsAffected(res sql.Result, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
