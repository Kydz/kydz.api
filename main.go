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
	k := Kouter.GetK()

	k.Get("/article", handlers.GetArticles)
	k.Post("/article", handlers.PostArticle)
	k.Get("/article/", handlers.GetArticle)
	k.Put("/article/", handlers.PutArticle)
	k.Delete("/article/", handlers.DelArticle)
	log.Println("server start")
	log.Fatal(http.ListenAndServe(":8088", k))
}
