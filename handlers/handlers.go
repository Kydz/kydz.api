package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Kydz/kydz.api/business"
	"log"
	"net/http"
	"strconv"
)

func ArticlesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Content-Type", "application/xml")
	switch r.Method {
	case http.MethodGet:
		articleDTO := &business.ArticleDTO{}
		page := getPage(r)
		perpage := getPerpage(r)
		rows := articleDTO.QueryList(page, perpage)
		jsonizeResponse(rows, w, r)
		break
	default:
		http.NotFound(w, r)
	}
}

func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Content-Type", "application/xml")
	param := r.URL.Path[len("/article/"):]
	id := convertToInt(param)
	log.Printf("ID ID ID %s, param %s", string(id), param)
	articleDTO := &business.ArticleDTO{}
	switch r.Method {
	case http.MethodGet:
		article := articleDTO.QuerySingle(id)
		jsonizeResponse(article, w, r)
		break
	case http.MethodPut:
		err := articleDTO.UpdateSingle(id, `content`, r.FormValue("content"))
		if err != nil {
			log.Fatal(err)
			errorResponse("Update article failed", w, r)
		}
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

func jsonizeResponse(data interface{}, w http.ResponseWriter, r *http.Request) {
	jsonResponse, err := json.Marshal(data)
	var response string
	if err != nil {
		log.Fatal(nil)
		errorResponse("Jsonize filed", w, r)
	} else {
		response = string(jsonResponse)
		sendResponse(response, w, r)
	}

}

func errorResponse(message string, w http.ResponseWriter, r *http.Request) {
	var response = `{"error": true, "message":` + message + `}`
	sendResponse(response, w, r)
}

func sendResponse(response string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprint(w, string(response))
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
