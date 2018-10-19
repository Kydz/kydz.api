package business

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)



var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("mysql", "root:1989222@/kydz")
	if err != nil {
		log.Print("Error: Opening [kydz] failed:")
		panic(err)
	}
}

type Article struct {
	Id      int
	Title   string
	Brief   string
	Content string
}
