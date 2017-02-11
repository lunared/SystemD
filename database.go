package main

//Note, this is currently unused, until I get around to adding database functionality this file will not be cleaned up

import (
	_ "database/sql"

	"github.com/jmoiron/sqlx"
	//_ "github.com/lib/pq"
	_ "log"
)

type Count struct {
	CountVal int `db:"count_val"`
}

//TODO: set up database to use json config file
func getUpdateCount() (int, error) {
	db, err := sqlx.Connect("postgres", "user=xxx password=xxx dbname=xxx sslmode=disable")
	if err != nil {
		return 0, err
	}

	count := Count{}

	db.Get(&count, "SELECT * FROM irc.count ORDER BY count_val")
	count.CountVal += 1

	tx := db.MustBegin()
	tx.MustExec("UPDATE irc.count set count_val = $1;", count.CountVal)
	tx.Commit()

	return count.CountVal, nil
}
