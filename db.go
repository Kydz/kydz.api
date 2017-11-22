package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type ArticleDTO struct{}

func (ad *ArticleDTO) QueryList(page int, perpage int) []Article {
	var list = make([]Article, perpage)
	db := getDB()
	offset := page * perpage
	limit := perpage
	queryString := `SELECT id, title, brief, content FROM articles WHERE active = 1 limit ?, ?`
	statement, err := db.Prepare(queryString)
	if err != nil {
		fmt.Println("Error: prepare [article sql] failed:")
		panic(err)
	}
	rows, err := statement.Query(offset, limit)
	if err != nil {
		fmt.Println("Error: query [article] failed:")
		panic(err)
	}
	defer rows.Close()
	var index = 0
	for rows.Next() {
		var id int
		var title string
		var brief string
		var content string
		rows.Scan(&id, &title, &brief, &content)
		list[index] = Article{Id: id, Title: title, Brief: brief, Content: content}
		index++
	}
	return list
}

func getDB() *sql.DB {
	db, err := sql.Open("mysql", "")
	if err != nil {
		fmt.Println("Error: Opening [kydz] failed:")
		panic(err)
	}
	return db
}

type Article struct {
	Id      int
	Title   string
	Brief   string
	Content string
}
