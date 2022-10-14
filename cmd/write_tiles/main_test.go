package main

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/mattn/go-sqlite3"
)

func db(file string) *sql.DB {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		panic(err)
	}
	return db
}

func Test_one(t *testing.T) {
	dir, _ := os.Getwd()
	log.Printf("dir: %s", dir)
	dbc := db("../data.output.sqlite")
	defer dbc.Close()
	_ = dbc
}

func Test_two(t *testing.T) {
	c1, c2 := Usbbox()
	spew.Dump(Usbbox())
	spew.Dump(RoundBlock(c1, c2))

	log.Printf("%d", ((c2.X-c1.X)*(c2.Y-c1.Y))/256.0)

}

func Test_usa(t *testing.T) {
	dbc := db("../data/output.mbtiles")
	e := WriteUsa(dbc, "../data/flat")
	if e != nil {
		panic(e)
	}
}
