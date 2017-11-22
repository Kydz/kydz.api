package db

import (
	"database/sql"
	"log"

	"github.com/davecgh/go-spew/spew"
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
		log.Print("Error: prepare [articles sql] failed:")
		panic(err)
	}
	rows, err := statement.Query(offset, limit)
	defer rows.Close()
	if err != nil {
		log.Print("Error: query [article] failed:")
		panic(err)
	}
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

func (ad *ArticleDTO) QuerySingle(queryId int) (article Article) {
	db := getDB()
	queryString :=`SELECT id, title, brief, content FROM articles WHERE id = ? AND active = 1`
	statement, err := db.Prepare(queryString)
	if err != nil {
		log.Print("Error: prepare [article sql] failed:")
		panic(err)
	}
	row, err := statement.QueryRow(queryId)
	defer row.Close()
	if err != nil {
		log.Printf("Error: query single %s failed:", string(queryId))
		panic(err)
	}
	var id int
	var title string
	var brief string
	var content string
	row.Scan(&id, &title, &brief, &content)
	article = Article{Id: id, Title: title, Brief: brief, Content: content}
	return article
}

func getDB() *sql.DB {
	db, err := sql.Open("mysql", "root@/kydz")
	if err != nil {
		log.Print("Error: Opening [kydz] failed:")
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
