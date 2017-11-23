package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Kydz/kydz.api"
)

func articlesHanlder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Content-Type", "application/xml")
	switch r.Method {
	case http.MethodGet:
		articleDTO := &db.ArticleDTO{}
		page := getPage(r)
		perpage := getPerpage(r)
		rows := articleDTO.QueryList(page, perpage)
		list, err := json.Marshal(rows)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprint(w, string(list))
		break
	default:
		http.NotFound(w, r)
	}
}

func articleHanlder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Content-Type", "application/xml")
	switch r.Method {
	case http.MethodGet:
		articleDTO := &db.ArticleDTO{}
		param := r.URL.Path[len("/article/"):]
		id := convertToInt(param)
		article := articleDTO.QuerySingle(id)
		row, err := json.Marshal(article)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprint(w, string(row))
		break
	default:
		http.NotFound(w, r)
	}
}

func getPage(r *http.Request) int {
	value := getFromForm(r, "page", "0")
	page := convertToInt(value)
	return page
}

func getPerpage(r *http.Request) int {
	value := getFromForm(r, "perpage", "20")
	perpage := convertToInt(value)
	return perpage
}

func convertToInt(value string) int {
	converted, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal(err)
	}
	return converted
}

func getFromForm(r *http.Request, key string, defaultValue string) string {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	value := r.Form.Get(key)
	if len(value) == 0 {
		value = defaultValue
	}
	return value
}

func main() {
	http.HandleFunc("/article", articlesHanlder)
	http.HandleFunc("/article/", articleHanlder)
	err := http.ListenAndServe(":8088", nil)
	log.Fatal(err)
}
