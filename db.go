package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/russross/blackfriday.v2"
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
		var content []byte
		rows.Scan(&id, &title, &brief, &content)
		if id > 0 {
			contentHTML := parseMarkdonw(content)
			list[index] = Article{Id: id, Title: title, Brief: brief, Content: contentHTML}
			index++
		} else {
			continue
		}
	}
	db.Close()
	return list
}

func (ad *ArticleDTO) QuerySingle(queryId int) (article Article) {
	db := getDB()
	queryString := `SELECT id, title, brief, content FROM articles WHERE id = ? AND active = 1`
	statement, err := db.Prepare(queryString)
	if err != nil {
		log.Print("Error: prepare [article sql] failed:")
		panic(err)
	}
	row := statement.QueryRow(queryId)
	if err != nil {
		log.Printf("Error: query single %s failed:", string(queryId))
		panic(err)
	}
	var id int
	var title string
	var brief string
	var content []byte
	row.Scan(&id, &title, &brief, &content)
	contentHTML := parseMarkdonw(content)
	article = Article{Id: id, Title: title, Brief: brief, Content: contentHTML}
	db.Close()
	return article
}

func getDB() *sql.DB {
	db, err := sql.Open("mysql", "root:1989222@/kydz")
	if err != nil {
		log.Print("Error: Opening [kydz] failed:")
		panic(err)
	}
	return db
}

func parseMarkdonw(input []byte) string {
	output := blackfriday.Run(input)
	html := string(output)
	return html
}

type Article struct {
	Id      int
	Title   string
	Brief   string
	Content string
}
