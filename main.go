package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/PuerkitoBio/goquery"
)

type Item struct {
	Service     string
	Description string
}

func main() {
	supportedLanguages := []string{"ar", "cn", "de", "en", "es", "fr", "id", "it", "jp", "ko", "pt", "ru", "th", "tr", "tw", "vi"}
	if len(os.Args) > 1 {
		arg := os.Args[1]
		items := getItems(arg)
		checkLanguageSupport(supportedLanguages, arg)
		output(items)
	} else {
	}
}

func getItems(lang string) []Item {
	baseUrl := "https://aws.amazon.com/"
	doc, err := goquery.NewDocument(baseUrl + lang)
	checkError(err)
	res := []Item{}
	items := doc.Find(".lb-content-item")
	items.Each(func(_ int, item *goquery.Selection) {
		name := item.Find("span").Text()
		description := item.Find("cite").Text()
		path, _ := item.Find("a").Attr("href")
		path = baseUrl + path[1:]
		name = "[" + name + "](" + path + ")"
		res = append(res, Item{Service: name, Description: description})
	})
	return res
}

func output(items []Item) {
	test_relpath := "[test github relpath](./go.mod)\n\n"
	header := "| | |\n| - | - |\n"
	content := ""
	for _, item := range items {
		content = content + "| " + item.Service + " | " + item.Description + " |\n"
	}
	res := test_relpath + header + content
	fmt.Println(res)
	err := ioutil.WriteFile("./README.md", []byte(res), 0644)
	checkError(err)
}

func checkLanguageSupport(langs []string, lang string) {
	is := false
	for _, v := range langs {
		if lang == v {
			is = true
		}
	}
	if is == false {
		fmt.Println(lang + " language is not supported\nsupported languages:")
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
