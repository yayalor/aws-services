package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type TableRow struct {
	Service     string
	Description string
}

type Options struct {
	FlagSet  *flag.FlagSet
	Output   string
	Template string
}

const defaultLanguage = "en"

var (
	supportedLanguages = []string{"ar", "cn", "de", "en", "es", "fr", "id", "it", "jp", "ko", "pt", "ru", "th", "tr", "tw", "vi"}
	single             = false
	options            Options
)

func init() {
	flag.CommandLine.Init("aws-services", flag.ExitOnError)
	options.FlagSet = flag.NewFlagSet("aws-services", flag.ExitOnError)
	options.FlagSet.StringVar(&options.Output, "o", "README.md", "output path on single")
	options.FlagSet.StringVar(&options.Template, "t", "", "template path on single")
}

func main() {
	if len(os.Args) > 1 {
		single = true
		flag.Parse()
		args := flag.Args()
		options.FlagSet.Parse(args[1:])
		arg := os.Args[1]
		if err := checkTemplate(); err != nil {
			log.Fatal(err)
		}
		if err := checkLanguageSupport(arg); err != nil {
			log.Fatal(err)
		}
		if err := output(arg, arg == defaultLanguage); err != nil {
			log.Fatal(err)
		}
	} else {
		for _, lang := range supportedLanguages {
			if err := output(lang, lang == defaultLanguage); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func checkTemplate() error {
	path := options.Template
	if path != "" {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return errors.New(path + ": no such file")
		}
	}
	return nil
}

func isExistPath(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func checkLanguageSupport(lang string) error {
	is := false
	for _, v := range supportedLanguages {
		if lang == v {
			is = true
		}
	}
	if is == false {
		errMsg := lang + " language is not supported\nsupported languages:"
		for _, v := range supportedLanguages {
			errMsg = v + "\n"
		}
		return errors.New(errMsg)
	}
	return nil
}

func getTableRowItems(lang string) ([]TableRow, error) {
	res := []TableRow{}
	baseUrl := "https://aws.amazon.com/"
	url := baseUrl
	if lang != "en" {
		url = baseUrl + lang + "/products"
	} else {
		url = baseUrl + "/products"
	}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return res, err
	}
	items := doc.Find(".lb-content-item")
	items.Each(func(_ int, item *goquery.Selection) {
		name := item.Find("span").Text()
		description := item.Find("cite").Text()
		path, _ := item.Find("a").Attr("href")
		if lang != "en" {
			path = baseUrl + path[1:]
		}
		name = "[" + name + "](" + path + ")"
		res = append(res, TableRow{Service: name, Description: description})
	})
	return res, err
}

func getNav(isDefaultOutput bool) string {
	res := ""
	for _, lang := range supportedLanguages {
		isDefaultLanguage := lang == defaultLanguage
		if isDefaultOutput {
			if !isDefaultLanguage {
				res = res + " | [" + lang + "](./languages/README." + lang + ".md)"
			}
		} else {
			if isDefaultLanguage {
				res = res + " | [" + lang + "](../README.md)"
			} else {
				res = res + " | [" + lang + "](./README." + lang + ".md)"
			}
		}
	}
	return res + " |\n"
}

func output(lang string, isDefault bool) error {
	var err error
	items, err := getTableRowItems(lang)
	if err != nil {
		return err
	}
	var nav string
	if !single {
		nav = getNav(isDefault)
	}
	tableHeader := "| Service | Description |\n| --- | --- |\n"
	tableContent := ""
	for _, item := range items {
		tableContent = tableContent + "| " + item.Service + " | " + item.Description + " |\n"
	}
	res := nav + "\n" + tableHeader + tableContent
	if _, err := os.Stat("./languages"); os.IsNotExist(err) {
		if err := os.Mkdir("./languages", 0755); err != nil {
			return err
		}
	}
	if isDefault || single {
		out := "./README.md"
		if single {
			out = options.Output
		}
		if options.Template != "" {
			type Template struct {
				Content string
				Date    string
			}
			var tp Template
			tp.Content = res
			tp.Date = time.Now().Format("2006/01/02")
			tpl, err := template.ParseFiles(options.Template)
			if err != nil {
				return err
			}
			nf, err := os.Create(out)
			if err != nil {
				return err
			}
			defer nf.Close()
			err = tpl.Execute(nf, tp)
			if err != nil {
				return err
			}
		} else {
			if err := ioutil.WriteFile(out, []byte(res), 0644); err != nil {
				return err
			}
		}
	} else {
		if err := ioutil.WriteFile("./languages/README."+lang+".md", []byte(res), 0644); err != nil {
			return err
		}
	}
	return err
}
