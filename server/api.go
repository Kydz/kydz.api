package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Kydz/kydz.api"
)

func articleHanlder(w http.ResponseWriter, r *http.Request) {
	urls := strings.Split(r.URL.Path, "?")
	switch r.Method {
	case http.MethodGet:
		path := urls[0][len("/article"):]
		log.Print(path)
		if len(path) == 0 {
			articleDTO := &db.ArticleDTO{}
			page, _ := strconv.Atoi(r.Form.Get("page"))
			perpage, _ := strconv.Atoi(r.Form.Get("perpage"))
			rows := articleDTO.QueryList(page, perpage)
			list, _ := json.Marshal(rows)
			fmt.Fprint(w, string(list))

		} else {
			// id := path[1:]
		}
		break
	case http.MethodPost:
		break

	}
}

func main() {
	http.HandleFunc("/article", articleHanlder)
	err := http.ListenAndServe(":8088", nil)
	log.Fatal(err)
}
