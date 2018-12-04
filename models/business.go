package models

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Kydz/kydz.api/Konfigurator"
	"github.com/Kydz/kydz.api/utils"
	"log"
	"time"

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

func Login(rp string) (string, error) {
	a := queryAdmin()
	h := sha256.New()
	h.Write([]byte(a.Salt + rp))
	p := fmt.Sprintf("%x", h.Sum(nil))
	if p == a.Password {
		return maskAdminToken(a.Token), nil
	}
	return "", errors.New("wrong password")
}

func CheckAdminToken(t string) bool {
	a := queryAdmin()
	return t == maskAdminToken(a.Token)
}

func InitAdmin() error {
	h := sha256.New()
	s := []byte(kon.AdminPass)
	h.Write(s)
	p := fmt.Sprintf("%x", h.Sum(nil))
	err := initAdmin(utils.GenerateRandomString(8), string(p))
	return err
}

func maskAdminToken(t string) string {
	h := sha256.New()
	n := time.Now()
	h.Write([]byte(t + utils.IntegerToString(n.Year()) + utils.IntegerToString(int(n.Month()))))
	t = fmt.Sprintf("%x", h.Sum(nil))
	return string(t)
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
