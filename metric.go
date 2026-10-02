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
	ID   int
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

func getMetricByID(db *sql.DB, id int) (*Metric, error) {
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

func getMetrics(db *sql.DB) ([]*Metric, error) {
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

func updateMetricByID(db *sql.DB, mu *MetricUpdate) (int64, error) {
	// do a check that metric exists first, for now we assume that it does
	sqlQuery := `UPDATE metrics SET name = $1 WHERE id = $2`
	res, err := db.Exec(sqlQuery, mu.Name, mu.ID)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, nil
	}
	return rowsAffected, nil
}

// solve the question with MetricValues. Do we need to delete values before the metric?
// are we working with values separately or we deleting them beforehand
func deleteMetricByID(db *sql.DB, id int) (int64, error) {
	sqlQuery := `DELETE FROM metrics WHERE id = $id`
	res, err := db.Exec(sqlQuery, id)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, nil
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

func createMetricValue(db *sql.DB, mv *MetricValue) (int64, error) {
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

// parsing time from the following format
// i have two metrics in mind currently, those are: weight and steps.
// both those metrics doesn't require exact second when written down
const timeLayout = "2006-01-02 15:04"

func parseTime(s string) (time.Time, error) {
	return time.ParseInLocation(timeLayout, s, time.Local)
}

func notReallyMain() {
	db, err := sql.Open("sqlite3", "test.db?_foreign_keys=on")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// think about making this naming shorter. add-value sounds clumsy
	// add-values sounds even more clumsy. hard to proceed here with those names.
	// the working name would be `add` and that's about it here.
	//
	// os.Args cannot have length 0 becuase it's always at least the program name
	if len(os.Args) == 1 {
		fmt.Print(`METH (METrics Health): stupid little program to track metrics:
	-find: find ID
	-create: create ID name
	-add-value: add-value ID value [created_at]
	-update: update ID new-name
	-delete: delete ID
	-list: list
`)
		return
	}

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
			met, err := getMetricByID(db, id)
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
				fmt.Printf("Following error has occured: %s\n", err.Error())
				return
			}

			fmt.Printf("Metric created, rows affected: %d\n", rowsAffected)
			return
		}
	} else if os.Args[1] == "add-value" {
		if len(os.Args[2:]) != 2 && len(os.Args[2:]) != 3 {
			fmt.Printf("usage: add-value ID value [created_at]")
			return
		}
		var mv *MetricValue
		if len(os.Args[2:]) == 2 {
			id, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Printf("Following error occured: %s\n", err.Error())
				return
			}

			value, err := strconv.ParseFloat(os.Args[3], 64)
			if err != nil {
				fmt.Printf("Following error occured: %s\n", err.Error())
				return
			}

			// here we can write a simple metric check first. if there is no metric with the following ID we freaking suck
			mv = &MetricValue{MetricID: id, Value: value, CreatedAt: time.Now()}
		}

		if len(os.Args[2:]) == 3 {
			id, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Printf("Following error occured: %s\n", err.Error())
				return
			}

			value, err := strconv.ParseFloat(os.Args[3], 64)
			if err != nil {
				fmt.Printf("Following error occured: %s\n", err.Error())
				return
			}

			createdAt, err := parseTime(os.Args[4])
			if err != nil {
				fmt.Printf("Following error occured: %s\n", err.Error())
				return
			}

			// here we can write a simple metric check first. if there is no metric with the following ID we freaking suck
			mv = &MetricValue{MetricID: id, Value: value, CreatedAt: createdAt}
		}

		// for now it's empty
		_, err := createMetricValue(db, mv)
		if err != nil {
			fmt.Printf("Following error occured: %s\n", err.Error())
			return
		}
		fmt.Printf("Metric value added nigga - metricID: %d, value: %.2f, createdAt: %s\n", mv.MetricID, mv.Value, mv.CreatedAt.Format(timeLayout))
		return
	} else if os.Args[1] == "update" {
		if len(os.Args[2:]) != 2 {
			fmt.Printf("usage: update ID new-name\n")
			return
		} else {
			id, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Printf("Following error has occured: %s\n", err.Error())
				return
			}
			name := os.Args[3]
			mu := &MetricUpdate{ID: id, Name: name}

			rowsAffected, err := updateMetricByID(db, mu)
			if err != nil {
				fmt.Printf("Following error has occured: %s\n", err.Error())
				return
			}

			fmt.Printf("Metric created, rows affected: %d\n", rowsAffected)
			return
		}
	} else if os.Args[1] == "list" {
		if len(os.Args[2:]) != 0 {
			fmt.Printf("usage: list\n")
		} else {
			mets, err := getMetrics(db)
			if err != nil {
				fmt.Printf("Following error has occured: %s\n", err.Error())
				return
			}

			if len(mets) > 5 {
				fmt.Printf("Following metrics exist (cut to first 5):\n")
				for i := range 5 {
					fmt.Printf("\tMetric, ID: %d, Name: %s\n", mets[i].ID, mets[i].Name)
				}
				return
			} else {
				fmt.Printf("Following metrics exist:\n")
				for _, m := range mets {
					fmt.Printf("\tMetric, ID: %d, Name: %s\n", m.ID, m.Name)
				}
				return
			}
		}
	} else if os.Args[1] == "delete" {
		if len(os.Args[2:]) != 1 {
			fmt.Printf("usage: delete ID\n")
			return
		} else {
			id, err := strconv.Atoi(os.Args[2])
			if err != nil {
				fmt.Printf("Following error has occured: %s\n", err.Error())
				return
			}
			rowsAffected, err := deleteMetricByID(db, id)
			if err != nil {
				fmt.Printf("Following error has occured: %s\n", err.Error())
				return
			}

			fmt.Printf("Metric deleted, rows affected: %d\n", rowsAffected)
			return
		}
	} else {
		fmt.Printf("Unknown command, use find or create")
		return
	}
}
