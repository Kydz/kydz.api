package models

import (
	"database/sql"
	"github.com/Kydz/kydz.api/utils"
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
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Brief   string `json:"brief"`
	Content string `json:"content"`
	Active  int    `json:"active"`
}

func QueryArticleList(offset int, limit int) []Article {
	var list []Article
	queryString := `SELECT id, title, brief FROM articles WHERE active = 1 limit ?, ?`
	log.Printf("db-excuted %d, %d, %s", offset, limit, queryString)
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
	for rows.Next() {
		var a Article
		rows.Scan(&a.Id, &a.Title, &a.Brief)
		list = append(list, a);
	}
	return list
}

func QueryArticleSingle(queryId int) (a Article) {
	queryString := `SELECT id, title, brief, content, active FROM articles WHERE id = ? AND active = 1`
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
	err = row.Scan(&a.Id, &a.Title, &a.Brief, &a.Content, &a.Active)
	if err != nil {
		panic(err)
	}
	return a
}

func UpdateArticleSingle(queryId int, article Article) (err error) {
	queryString := `UPDATE articles SET `
	if article.Brief != "" {
		queryString += `brief = "` + article.Brief + `", `
	}
	if article.Title != "" {
		queryString += `title = "` + article.Title + `", `
	}
	if article.Content != "" {
		queryString += `content = "` + article.Content + `", `
	}
	queryString += `active = ` + utils.IntegerToString(article.Active) + ` WHERE id = ` + utils.IntegerToString(queryId) + `;`
	log.Printf("Excuted SQL: %s", queryString)
	_, err = db.Exec(queryString)
	if err != nil {
		log.Printf("Error: update article %s failed:", string(queryId))
		panic(err)
	}
	return err
}
