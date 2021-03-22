package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	supportedLanguages := []string{"ar", "cn", "de", "en", "es", "fr", "id", "it", "jp", "ko", "pt", "ru", "th", "tr", "tw", "vi"}
	fmt.Println(supportedLanguages)

	lang := os.Args[1]
	baseUrl := "https://aws.amazon.com/"
	fmt.Println(baseUrl, lang)
	doc, err := goquery.NewDocument(baseUrl + lang)
	checkError(err)

	items := doc.Find(".lb-content-item")
	items.Each(func(_ int, item *goquery.Selection) {
		path, _ := item.Find("a").Attr("href")
		fmt.Println(path)
	})
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func typeof(i interface{}) interface{} {
	return reflect.TypeOf(i)
}
