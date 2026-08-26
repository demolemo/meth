package sqlite

import (
	"database/sql"
)

type MetricService struct {
	db *sql.DB
}
