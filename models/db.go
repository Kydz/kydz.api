package models

import (
	"database/sql"
	"github.com/Kydz/kydz.api/Konfigurator"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var kon *Konfigurator.Kon

func init() {
	kon = Konfigurator.GetKon()
}

func InitDB() {
	var err error
	db, err = sql.Open("mysql", kon.Db.Account + ":" + kon.Db.Pass + "@/" + kon.Db.Schema)
	if err != nil {
		log.Print("Error: Opening [kydz] failed:")
		panic(err)
	}
}

func logQuery(qs string, args ... interface{}) {
	ls := `[execute DB query]:` + qs
	if len(args) > 0 {
		log.Printf(ls, args)
	} else {
		log.Println(ls)
	}
}

func GetArticleList(offset int, limit int) (list ArticleList) {
	list.Total = queryArticleTotal()
	if list.Total < offset {
		list.Rows = make([]Article, 0)
	} else {
		list.Rows = queryArticleList(offset, limit)
	}
	return list
}

func GetArticle(id int, hit bool) (a Article) {
	return queryArticleSingle(id, hit)
}

func PutArticle(id int, a Article) error {
	a.Id = id
	return updateArticleSingle(a)
}

func PostArticle(a Article) (id int64, err error) {
	id, err = insertArticle(a)
	return id, err
}

func DelArticle(id int) (int64, error) {
	return deleteArticle(id)
}

type ArticleList struct {
	Total int       `json:"total"`
	Rows  []Article `json:"rows"`
}

type Article struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Brief   string `json:"brief"`
	Content string `json:"content"`
	Active  int    `json:"active"`
	Hit     int    `json:"hit"`
	Type    int    `json:"type"`
}
