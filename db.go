package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type ArticleDTO struct {}

func (ad *ArticleDTO) QueryList (page int, perpage int) []Article {
	var list = make([]Article, perpage)
	db := getDB()
	offset := page * perpage
	limit := offset + perpage
	queryString := `SELECT id, title, brief, content FORM articles WHERE active = 1 limit ?, ?`
	statement, _ := db.Prepare(queryString)
	fmt.Print(offset)
	rows, err := statement.Query(offset, limit)
	if err != nil {
		fmt.Printf("Error: query [article] failed - %s", err)
		panic(err)
	}
	defer rows.Close()
	var index int = 1
	for rows.Next() {
		var id int
		var title string
		var brief string
		var content []byte
		rows.Scan(&id, &title, &brief, &content)
		list[index] = Article{Id: id, Title: title, Brief: brief, Content: content}
	}
	return list
}

func getDB() *sql.DB {
	db, err := sql.Open("mysql", "root@/kydz")
	if err != nil {
		fmt.Printf("Error: Opening [kydz] failed - %s", err)
		panic(err)
	}
	return db
}

type Article struct {
	Id int
	Title string
	Brief string
	Content []byte
}
