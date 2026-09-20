package main

import (
	"database/sql"
	"fmt"
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

func selectValues(db *sql.DB, ID string) error {
	return nil
}

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	affectedRows, err := insertValues(db, "random-id", "jora-gay", 1488)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Sir, some rows were affected and we are moving freaking forward: %d", affectedRows)
}
