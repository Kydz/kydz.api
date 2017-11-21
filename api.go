package main

import (
	"fmt"
	"github.com/Kydz/kydz.db"
)

func main() {
	articleDTO := &db.ArticleDTO{}
	rows := articleDTO.QueryList(0, 10)
	fmt.Print(rows)
}