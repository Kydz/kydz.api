package main

import (
	"github.com/Kydz/kydz.api/Kouter"
	"github.com/Kydz/kydz.api/handlers"
	"github.com/Kydz/kydz.api/models"
	"log"
	"net/http"
)

func main() {
	models.InitDB()
	k := Kouter.NewK()

	k.Get("article", handlers.GetArticles)
	k.Post("article", handlers.PostArticle)
	k.Get("article/{:id}", handlers.GetArticle)
	k.Put("article/{:id}", handlers.PutArticle)
	k.Delete("article/{:id}", handlers.DelArticle)
	log.Println("server start")
	log.Fatal(http.ListenAndServe(":8088", k))
}
