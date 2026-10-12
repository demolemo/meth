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

// upsertValue is the one write path for a reading. (metric_id, created_at) is
// a reading's identity and metric_values_reading enforces it, so the same
// instant arriving twice lands on the same row.
//
// created_at is bound as UTC everywhere below: go-sqlite3 stores a time as
// text carrying its offset, so 07:00Z and 10:00+03:00 are one instant and two
// different index keys unless they are normalised first.
const upsertValue = `
	INSERT INTO metric_values(metric_id, value, created_at, updated_at)
	VALUES($1, $2, $3, $4)
	ON CONFLICT(metric_id, created_at) DO UPDATE
	SET value = excluded.value, updated_at = excluded.updated_at
`

type ValueService struct {
	db *sql.DB
}

func NewValueService(db *sql.DB) *ValueService {
	return &ValueService{db: db}
}

func (vs *ValueService) CreateMetricValue(mv *meth.MetricValue) (int64, error) {
	// we are worrying about corretly supplying the date outside of this function
	// this solution gives us a little bit more responsibility
	res, err := vs.db.Exec(upsertValue, mv.MetricID, mv.Value, mv.CreatedAt.UTC(), time.Now().UTC())
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
		WHERE id = $1
	`
	mv := &meth.MetricValue{}
	err := vs.db.QueryRow(sqlQuery, valueID).Scan(&mv.ID, &mv.MetricID, &mv.Value, &mv.CreatedAt, &mv.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("value %d not found", valueID)
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
	return scanValues(rows)
}

// GetValuesBetween returns every metric's values in [from, to), newest first.
// The caller joins metric_id to a name itself - with one table of readings and
// no grouping anywhere, this is what "what was recorded that night" is.
func (vs *ValueService) GetValuesBetween(from, to time.Time) ([]*meth.MetricValue, error) {
	sqlQuery := `
		SELECT id, metric_id, value, created_at, updated_at FROM metric_values
		WHERE created_at >= $1 AND created_at < $2
		ORDER BY created_at DESC, metric_id
	`
	rows, err := vs.db.Query(sqlQuery, from.UTC(), to.UTC())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanValues(rows)
}

// QUESTION: should we send a value update instead of the value itself?
func (vs *ValueService) UpdateMetricValue(mv *meth.MetricValue) (int64, error) {
	sqlQuery := `UPDATE metric_values SET value = $1, updated_at = $2 where id = $3 AND metric_id = $4`
	res, err := vs.db.Exec(sqlQuery, mv.Value, mv.UpdatedAt.UTC(), mv.ID, mv.MetricID)
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

// CreateMetricValues writes a batch in one transaction: re-importing an
// overlapping export changes nothing unless the number did.
//
// The lookup is not a second write path, it is where the three counts come
// from. sqlite's changes() reports 1 for a conflicting row even when the
// DO UPDATE was suppressed by a WHERE, so an upsert on its own cannot tell
// "already there" from "updated". Reading the stored value first both answers
// that and lets an unchanged reading be skipped entirely.
func (vs *ValueService) CreateMetricValues(mvs []*meth.MetricValue) (meth.ImportResult, error) {
	var res meth.ImportResult

	tx, err := vs.db.Begin()
	if err != nil {
		return res, err
	}
	defer tx.Rollback()

	lookup, err := tx.Prepare(`SELECT value FROM metric_values WHERE metric_id = $1 AND created_at = $2`)
	if err != nil {
		return res, err
	}
	defer lookup.Close()

	upsert, err := tx.Prepare(upsertValue)
	if err != nil {
		return res, err
	}
	defer upsert.Close()

	now := time.Now().UTC()
	for i, mv := range mvs {
		var stored sql.NullFloat64
		err := lookup.QueryRow(mv.MetricID, mv.CreatedAt.UTC()).Scan(&stored)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			res.Inserted++
		case err != nil:
			return res, fmt.Errorf("value %d of %d: %w", i+1, len(mvs), err)
		case stored.Valid && stored.Float64 == mv.Value:
			res.Unchanged++
			continue
		default:
			res.Updated++
		}

		if _, err := upsert.Exec(mv.MetricID, mv.Value, mv.CreatedAt.UTC(), now); err != nil {
			return res, fmt.Errorf("value %d of %d: %w", i+1, len(mvs), err)
		}
	}
	return res, tx.Commit()
}

func scanValues(rows *sql.Rows) ([]*meth.MetricValue, error) {
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
