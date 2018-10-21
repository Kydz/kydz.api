package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Kydz/kydz.api/models"
	"github.com/Kydz/kydz.api/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func ArticlesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Content-Type", "application/xml")
	switch r.Method {
	case http.MethodGet:
		articleDTO := &models.ArticleDTO{}
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
	id := utils.StringToInteger(param)
	articleDTO := &models.ArticleDTO{}
	switch r.Method {
	case http.MethodGet:
		article := articleDTO.QuerySingle(id)
		jsonizeResponse(article, w, r)
		break
	case http.MethodPut:
		bodyBytes := readJsonRequestBody(r.Body)
		if bodyBytes == nil {
			log.Fatal("Empty request body")
			errorResponse("Empty request body", w, r)
		}
		article, err := parseJsonToArticle(bodyBytes)
		if err != nil {
			log.Fatal(err)
			errorResponse("Parse Json failed, got: " + string(bodyBytes), w, r)
		}
		err = articleDTO.UpdateSingle(id, article)
		if err != nil {
			log.Fatal(err)
			errorResponse("Update article failed", w, r)
		} else {
			normalResponse("{\"success\": true}", w, r)
		}
		break
	case http.MethodOptions:
		corsResponse(w, r)
		break
	default:
		http.NotFound(w, r)
	}
}

func getPage(r *http.Request) int {
	value := getFromForm(r, "page", "0")
	page := utils.StringToInteger(value)
	return page
}

func getPerpage(r *http.Request) int {
	value := getFromForm(r, "perpage", "20")
	perpage := utils.StringToInteger(value)
	return perpage
}

func jsonizeResponse(data interface{}, w http.ResponseWriter, r *http.Request) {
	jsonResponse, err := json.Marshal(data)
	var response string
	if err != nil {
		log.Fatal(err)
		errorResponse("Jsonize filed", w, r)
	} else {
		response = string(jsonResponse)
		normalResponse(response, w, r)
	}
}

func parseJsonToArticle(jsonBytes []byte) (article models.Article, err error) {
	err = json.Unmarshal(jsonBytes, &article)
	return article, err
}

func errorResponse(message string, w http.ResponseWriter, r *http.Request) {
	var response = `{"error": true, "message":` + message + `}`
	normalResponse(response, w, r)
}

func normalResponse(response string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprint(w, string(response))
}

func corsResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func readJsonRequestBody(reader io.Reader) (bytes []byte) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return bytes
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
