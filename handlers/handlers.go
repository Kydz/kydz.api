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
	switch r.Method {
	case http.MethodGet:
		qp := r.URL.Query()
		offset := 0
		limit := 20
		if qp.Get("o") != "" {
			offset = utils.StringToInteger(qp.Get("o"))
		}
		if qp.Get("l") != "" {
			limit = utils.StringToInteger(qp.Get("l"))
		}
		rows := models.GetArticleList(offset, limit)
		jsonizeResponse(rows, w, r)
		break
	case http.MethodPost:
		a, err := getArticleFromRequestBody(r.Body)
		if err != nil {
			log.Fatal(err)
			errorResponse("Parse Json failed", w, r)
		}
		id, err := models.PostArticle(a)
		if err != nil {
			log.Fatal(err)
			errorResponse("Add article failed", w, r)
		} else {
			normalResponse("{\"success\": true, \"id\": "+utils.IntegerToString(int(id))+"}", w, r)
		}
		break
	case http.MethodOptions:
		corsResponse(w, r)
		break
	default:
		http.NotFound(w, r)
	}
}

func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	pm := r.URL.Path[len("/article/"):]
	id := utils.StringToInteger(pm)
	qp := r.URL.Query()
	log.Printf("query params get: %+v", qp)
	switch r.Method {
	case http.MethodGet:
		hit := true
		if qp.Get("h") != "" {
			hit = false
		}
		article := models.GetArticle(id, hit)
		jsonizeResponse(article, w, r)
		break
	case http.MethodPut:
		a, err := getArticleFromRequestBody(r.Body)
		if err != nil {
			log.Fatal(err)
			errorResponse("Parse Json failed", w, r)
		}
		err = models.PutArticle(id, a)
		if err != nil {
			log.Fatal(err)
			errorResponse("Update article failed", w, r)
		} else {
			normalResponse("{\"success\": true}", w, r)
		}
		break
	case http.MethodDelete:
		rows, err := models.DelArticle(id)
		if err != nil {
			log.Fatal(err)
			errorResponse("Parse Json failed", w, r)
		} else {
			normalResponse("{\"success\": true, \"affectedRows\": "+utils.IntegerToString(int(rows))+"}", w, r)
		}
		break
	case http.MethodOptions:
		corsResponse(w, r)
		break
	default:
		http.NotFound(w, r)
	}
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

func errorResponse(message string, w http.ResponseWriter, r *http.Request) {
	var response = `{"error": true, "message":` + message + `}`
	normalResponse(response, w, r)
}

func normalResponse(response string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(response))
}

func corsResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
	w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func getArticleFromRequestBody(reader io.Reader) (a models.Article, err error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bytes, &a)
	return a, err
}
