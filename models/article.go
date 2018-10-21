package models

import (
	"github.com/Kydz/kydz.api/utils"
	"log"
)

type ArticleDTO struct{}

func (ad *ArticleDTO) QueryList(page int, perpage int) []Article {
	var list = make([]Article, perpage)
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
			list[index] = Article{Id: id, Title: title, Brief: brief, Content: string(content)}
			index++
		} else {
			continue
		}
	}
	return list
}

func (ad *ArticleDTO) QuerySingle(queryId int) (article Article) {
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
	var id int
	var title string
	var brief string
	var content []byte
	var active int
	err = row.Scan(&id, &title, &brief, &content, &active)
	if err != nil {
		panic(err)
	}
	article = Article{Id: id, Title: title, Brief: brief, Content: string(content), Active: active}
	return article
}

func (ad *ArticleDTO) UpdateSingle(queryId int, article Article) (err error) {
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
