package main

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	doc, err := goquery.NewDocument("https://gobyexample.com/")
	checkErr(err)
	doc.Find("ul a").Each(func(_ int, s *goquery.Selection) {
		path, _ := s.Attr("href")
		con, err := goquery.NewDocument("https://gobyexample.com/" + path)
		checkErr(err)
		fmt.Println(con)

		con.Find("ul a").Each(func(_ int, s *goquery.Selection) {
		})
	})
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
