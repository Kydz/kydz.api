package handlers

import (
	"encoding/json"
	"github.com/Kydz/kydz.api/Kouter"
	"github.com/Kydz/kydz.api/models"
	"github.com/Kydz/kydz.api/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func GetArticles(w http.ResponseWriter, r *http.Request) {
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
	jsonizeResponse(rows, w)
}

func PostArticle(w http.ResponseWriter, r *http.Request) {
	a, err := getArticleFromRequestBody(r.Body)
	if err != nil {
		log.Fatal(err)
		errorResponse("Parse Json failed", w)
	}
	id, err := models.PostArticle(a)
	if err != nil {
		log.Fatal(err)
		errorResponse("Add article failed", w)
	} else {
		normalResponse("{\"success\": true, \"id\": "+utils.IntegerToString(int(id))+"}", w)
	}
}

func GetArticle(w http.ResponseWriter, r *http.Request) {
	id := utils.StringToInteger(Kouter.GetCurrentRoute().Params["id"])
	qp := r.URL.Query()
	hit := true
	if qp.Get("h") != "" {
		hit = false
	}
	article := models.GetArticle(id, hit)
	jsonizeResponse(article, w)
}


func PutArticle(w http.ResponseWriter, r *http.Request) {
	id := utils.StringToInteger(Kouter.GetCurrentRoute().Params["id"])
	a, err := getArticleFromRequestBody(r.Body)
	if err != nil {
		log.Fatal(err)
		errorResponse("Parse Json failed", w)
	}
	err = models.PutArticle(id, a)
	if err != nil {
		log.Fatal(err)
		errorResponse("Update article failed", w)
	} else {
		normalResponse("{\"success\": true}", w)
	}
}

func DelArticle(w http.ResponseWriter, r *http.Request) {
	id := utils.StringToInteger(Kouter.GetCurrentRoute().Params["id"])
	rows, err := models.DelArticle(id)
	if err != nil {
		log.Fatal(err)
		errorResponse("Parse Json failed", w)
	} else {
		normalResponse("{\"success\": true, \"affectedRows\": "+utils.IntegerToString(int(rows))+"}", w)
	}
}


func getArticleFromRequestBody(reader io.Reader) (a models.Article, err error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bytes, &a)
	return a, err
}
