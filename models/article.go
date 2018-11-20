package models

import (
	"log"
	"reflect"
	"strings"
)

func queryArticleList(offset int, limit int) []Article {
	var list []Article
	qs := `SELECT id, title, brief FROM articles WHERE active = 1 limit ?, ?`
	logQuery(qs+"; with: %d, %d, %s", offset, limit, qs)
	stmt, err := db.Prepare(qs)
	if err != nil {
		log.Print("Error: prepare [articles sql] failed:")
		panic(err)
	}
	rows, err := stmt.Query(offset, limit)
	if err != nil {
		log.Print("Error: query [article] failed:")
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var a Article
		rows.Scan(&a.Id, &a.Title, &a.Brief)
		list = append(list, a)
	}
	return list
}

func queryArticleTotal() (total int) {
	qs := `SELECT count(id) as total FROM articles WHERE active = 1`
	logQuery(qs)
	stmt, err := db.Prepare(qs)
	if err != nil {
		log.Print("Error: query article total failed")
		panic(err)
	}
	row := stmt.QueryRow()
	err = row.Scan(&total)
	if err != nil {
		panic(err)
	}
	return total
}

func queryArticleSingle(queryId int, updateHit bool) (a Article) {
	qs := `SELECT id, title, brief, content, active FROM articles WHERE id = ? AND active = 1`
	logQuery(qs, queryId)
	stmt, err := db.Prepare(qs)
	if err != nil {
		log.Print("Error: prepare [article sql] failed:")
		panic(err)
	}
	row := stmt.QueryRow(queryId)
	err = row.Scan(&a.Id, &a.Title, &a.Brief, &a.Content, &a.Active)
	if updateHit {
		addArticleHit(a.Id)
	}
	if err != nil {
		panic(err)
	}
	return a
}

func addArticleHit(id int) {
	qs := `UPDATE articles SET hit = hit + 1 WHERE id = ?`
	logQuery(qs + " %d", id)
	_, err := db.Exec(qs, id)
	if err != nil {
		log.Print("Error: execute update hit query failed")
		log.Print(err)
	}
}

func updateArticleSingle(a Article) (err error) {
	fields, values := getORM(a)
	qs := `UPDATE articles SET `
	qs += strings.Join(fields, " = ?,")
	qs += ` = ? WHERE id = ?`
	values = append(values, a.Id)
	logQuery(qs+" %+v", values)
	_, err = db.Exec(qs, values...)
	if err != nil {
		log.Printf("Error: update article %+v failed:", string(a.Id))
		panic(err)
	}
	return err
}

func insertArticle(a Article) (id int64, err error) {
	fields, values := getORM(a)
	qs := "INSERT INTO articles ("
	qs += strings.Join(fields, ", ")
	qs += ") VALUES ("
	for i := 1; i < len(fields); i ++ {
		qs += "? ,"
	}
	qs += "? );"
	logQuery(qs + "%+v", values)
	r, err := db.Exec(qs, values...)
	if err != nil {
		log.Printf("Error: insert into article failed, values: %+v", a)
		panic(err)
	}
	id, err = r.LastInsertId()
	if err != nil {
		log.Println("Error: execute insert failed")
		panic(err)
	}
	return id, err
}

func getORM(a Article) ([]string, []interface{}) {
	a.Type = 1
	elA := reflect.ValueOf(&a).Elem()
	typesA := elA.Type()
	lenA := elA.NumField()
	fields := make([]string, lenA)
	values := make([]interface{}, lenA)
	log.Print("in")
	log.Printf("%+v", a)
	j := 0
	for i := 0; i < lenA; i++ {
		value := elA.Field(i).Interface()
		if value != "" || value != nil {
			fields[j] = typesA.Field(i).Name
			values[j] = value
			j++
		}
	}
	return fields, values
}

func deleteArticle(id int) (int64, error) {
	qs := `DELETE FROM articles WHERE id = ?`
	res, err := db.Exec(qs, id)
	if err != nil {
		log.Print("Error: delete article failed")
		log.Print(err)
	}
	return res.RowsAffected()
}