package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// how does the simplest version of the program works?
// we can add values
// and edit values
// and view values
// and that's about it
// there are not 1000 of abstractions that fuck us over and make our life living hell?
// the issue with the approach of writing everything correctly are actually two
// 1. it's not fucking fun at all. i'm always thinking about how to make things right instead of how to make things work
// 2. we have no working copy of the program till very late into writing process
func insertValues(db *sql.DB, ID string, name string, salary int64) (int64, error) {
	sqlQuery := `INSERT INTO emp VALUES($1, $2, $3)`
	res, err := db.Exec(sqlQuery, ID, name, salary)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

type Metric struct {
	ID   int
	Name string
}

type MetricUpdate struct {
	Name string
}

func createMetric(db *sql.DB, m *Metric) (int64, error) {
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

func findMetricByID(db *sql.DB, id int) (*Metric, error) {
	sqlQuery := `
		SELECT * FROM metrics WHERE id = $id
	`
	rows, err := db.Query(sqlQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var name string

	// what would this do if there is no rows?
	for rows.Next() {
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		} else {
			return &Metric{ID: id, Name: name}, nil
		}
	}
	return nil, nil
}

func updateMetric(db *sql.DB, m *MetricUpdate) (int64, error) {
	return 0, nil
}

type MetricValue struct {
	ID       int
	MetricID int
	Value    float64
	Created  time.Time
	Updated  time.Time
}

func main() {
	db, err := sql.Open("sqlite3", "test.db?_foreign_keys=on")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if os.Args[1] == "find" {
		if len(os.Args[2:]) != 1 {
			fmt.Printf("usage: find ID\n")
			return
		} else {
			id, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Printf("Following error has occured: %s\n", err.Error())
				return
			}
			met, err := findMetricByID(db, id)
			if err != nil {
				fmt.Printf("Following error has occured: %s\n", err.Error())
				return
			}
			fmt.Printf("Metric, ID: %d, Name: %s\n", met.ID, met.Name)
			return
		}
	} else if os.Args[1] == "create" {
		if len(os.Args[2:]) != 2 {
			fmt.Printf("usage: create ID name\n")
			return
		} else {
			id, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Printf("Following error has occured: %s\n", err.Error())
				return
			}
			name := os.Args[3]
			met := &Metric{ID: id, Name: name}

			rowsAffected, err := createMetric(db, met)
			if err != nil {
				fmt.Printf("Follwing error has occured: %s\n", err.Error())
				return
			}

			fmt.Printf("Metric created, rows affected: %d\n", rowsAffected)
			return
		}
	} else {
		fmt.Printf("Unknown command, use find or create")
		return
	}
}
