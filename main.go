package main

import (
	"github.com/Kydz/kydz.api/business"
	"github.com/Kydz/kydz.api/handlers"
	"log"
	"net/http"
)

func main() {
	business.InitDB()
	http.HandleFunc("/", http.NotFound)
	http.HandleFunc("/article", handlers.ArticlesHandler)
	http.HandleFunc("/article/", handlers.ArticleHandler)
	log.Println("business start")
	err := http.ListenAndServe(":8088", nil)
	log.Fatal(err)
}
