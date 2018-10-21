package main

import (
	"github.com/Kydz/kydz.api/handlers"
	"github.com/Kydz/kydz.api/models"
	"log"
	"net/http"
)

func main() {
	models.InitDB()
	http.HandleFunc("/", http.NotFound)
	http.HandleFunc("/article", handlers.ArticlesHandler)
	http.HandleFunc("/article/", handlers.ArticleHandler)
	log.Println("server start")
	err := http.ListenAndServe(":8088", nil)
	log.Fatal(err)
}
