package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/PuerkitoBio/goquery"
)

type Item struct {
	Name        string
	Description string
	Path        string
}

func main() {
	arg := os.Args[1]
	supportedLanguages := []string{"ar", "cn", "de", "en", "es", "fr", "id", "it", "jp", "ko", "pt", "ru", "th", "tr", "tw", "vi"}
	if arg == "all" {
		fmt.Println("all")
	} else {
		items := getItems(arg)
		checkLanguageSurpport(supportedLanguages, arg)
		fmt.Println(items)
	}
}

func getItems(lang string) []Item {
	baseUrl := "https://aws.amazon.com/"
	fmt.Println(baseUrl, lang)
	doc, err := goquery.NewDocument(baseUrl + lang)
	checkError(err)
	res := []Item{}
	items := doc.Find(".lb-content-item")
	items.Each(func(_ int, item *goquery.Selection) {
		name := item.Find("span").Text()
		description := item.Find("cite").Text()
		path, _ := item.Find("a").Attr("href")
		res = append(res, Item{Name: name, Description: description, Path: path})
	})
	return res
}

func checkLanguageSurpport(langs []string, lang string) {
	is := false
	for _, v := range langs {
		if lang == v {
			is = true
		}
	}
	if is == false {
		fmt.Println(lang + " language is not surpported\nsurpported languages:")
		for _, v := range langs {
			fmt.Println(v)
		}
		os.Exit(0)
	}
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
