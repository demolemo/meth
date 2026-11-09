package main

import (
	"flag"
	"log"

	"github.com/demolemo/meth/config"
	"github.com/demolemo/meth/http"
	"github.com/demolemo/meth/sqlite"
)

func main() {
	dbFlag := flag.String("db", "", "sqlite file to use, overrides the config file")
	flag.Parse()

	path := *dbFlag
	if path == "" {
		cfg, err := config.Load()
		if err != nil {
			log.Fatal(err)
		}
		path = cfg.DB
	}

	db, err := sqlite.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Printf("db: %s", path)
	log.Fatal(http.NewServer(db).Serve())
}
