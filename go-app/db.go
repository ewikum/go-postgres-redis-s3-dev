package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"strings"
	"net/http"
	"database/sql"
	_"github.com/lib/pq"
)

var (
	db *sql.DB
)

type dbHandler struct{}

func (h dbHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "DB Page\n\n")

	err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	addnew:=strings.TrimPrefix(r.URL.Path, "/db/add/")
	if addnew!="" {
		insertOne(addnew)	
		fmt.Fprintf(w, "New item inserted. \n\n")
	}else {
		fmt.Fprintf(w, "Access /db/add/{some text} to insert new item. \n\n")
	}

	fmt.Fprintf(w, "Items from dummy table...\n\n")

	items:=listitems()
	for c,obkey  := range items {
		fmt.Fprintln(w,c+1,". ",obkey)
	}
}

func initDB() error {
	connInfo := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_ENV_POSTGRES_USER"),
		os.Getenv("DB_ENV_POSTGRES_DBNAME"),
		os.Getenv("DB_ENV_POSTGRES_PASSWORD"),
		os.Getenv("DB_PORT_5432_TCP_ADDR"),
		os.Getenv("DB_PORT_5432_TCP_PORT"),
	)

	var err error
	db, err = sql.Open("postgres", connInfo)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(i) * time.Second)

		if err = db.Ping(); err == nil {
			break
		}
		log.Println(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(
		`create table if not exists dummytable (
			id SERIAL,
			dummytext TEXT NOT NULL,
			CONSTRAINT dummytable_pkey PRIMARY KEY (id)
		)`)

	return err
}

func listitems() []string {

	rows, err := db.Query("select dummytext from dummytable")
	if err != nil {
		log.Fatal(err)
	}

	var items []string
    for rows.Next() {
		var dummytext string
        err := rows.Scan(&dummytext)
        if err != nil {
            log.Fatal(err)
        }
        items = append(items, dummytext)
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
	}

	
	return items
}

func insertOne(text string){
	r, err := db.Exec("insert into dummytable(dummytext) values('" + text + "')")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(r)
}