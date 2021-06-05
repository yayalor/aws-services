package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Item struct {
	Service     string
	Description string
}

func main() {
	defaultLanguage := "en"
	supportedLanguages := []string{"ar", "cn", "de", "en", "es", "fr", "id", "it", "jp", "ko", "pt", "ru", "th", "tr", "tw", "vi"}
	if len(os.Args) > 1 {
		arg := os.Args[1]
		items, err := getItems(arg)
		if err != nil {
			log.Fatal(err)
		}
		checkLanguageSupport(supportedLanguages, arg)
		if err := outputItem(items, supportedLanguages, arg, arg == defaultLanguage); err != nil {
			log.Fatal(err)
		}
	} else {
		for _, lang := range supportedLanguages {
			items, err := getItems(lang)
			if err != nil {
				log.Fatal(err)
			}
			checkLanguageSupport(supportedLanguages, lang)
			if err := outputItem(items, supportedLanguages, lang, lang == defaultLanguage); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func checkLanguageSupport(langs []string, lang string) {
	is := false
	for _, v := range langs {
		if lang == v {
			is = true
		}
	}
	if is == false {
		msg := strings.Join([]string{lang, " language is not supported\nsupported languages:"}, "")
		for _, v := range langs {
			msg = strings.Join([]string{v, "\n"}, "")
		}
		log.Fatalf(msg)
	}
}

func getItems(lang string) ([]Item, error) {
	res := []Item{}
	baseUrl := "https://aws.amazon.com/"
	url := strings.Join([]string{baseUrl, lang}, "")
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return res, err
	}
	items := doc.Find(".lb-content-item")
	items.Each(func(_ int, item *goquery.Selection) {
		name := item.Find("span").Text()
		description := item.Find("cite").Text()
		path, _ := item.Find("a").Attr("href")
		path = strings.Join([]string{baseUrl, path[1:]}, "")
		name = strings.Join([]string{"[", name, "](", path, ")"}, "")
		res = append(res, Item{Service: name, Description: description})
	})
	return res, err
}

func getNavs(langs []string, isDef bool) string {
	res := ""
	for _, lang := range langs {
		if isDef {
			res = strings.Join([]string{res, " | [", lang, "](./languages/README.", lang, ".md)"}, "")
		} else {
			res = strings.Join([]string{res, " | [", lang, "](./README.", lang, ".md)"}, "")
		}
	}
	return strings.Join([]string{res, " |\n"}, "")
}

func outputItem(items []Item, langs []string, lang string, isDef bool) error {
	var err error
	header := "| Service | Description |\n| - | - |\n"
	content := ""
	for _, item := range items {
		content = strings.Join([]string{content, "| ", item.Service, " | ", item.Description, " |\n"}, "")
	}
	navs := getNavs(langs, isDef)
	res := strings.Join([]string{navs, header, content}, "")
	if _, err := os.Stat("./languages"); os.IsNotExist(err) {
		if err := os.Mkdir("./languages", 0755); err != nil {
			return err
		}
	}
	if isDef {
		if err := WriteFile("./README.md", []byte(res), 0644); err != nil {
			return err
		}
	}
	if err := WriteFile(strings.Join([]string{"./languages/README.", lang, ".md"}, ""), []byte(res), 0644); err != nil {
		return err
	}
	return err
}

func WriteFile(path string, data []byte, perm fs.FileMode) error {
	err := ioutil.WriteFile(path, data, perm)
	return err
}
