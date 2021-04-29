package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func loadDocument(pitcherUrl string) *goquery.Document {
	res, err := http.Get(pitcherUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func names(doc *goquery.Document) []string {
	names := make([]string, 0)
	doc.Find(".item-title").Each(func(i int, s *goquery.Selection) {
		names = append(names, s.Text())
	})

	return names
}

func main() {
	pitcherUrl := "https://www.yakult-swallows.co.jp/players/category/pitcher"
	catcherUrl := "https://www.yakult-swallows.co.jp/players/category/catcher"
	infielderUrl := "https://www.yakult-swallows.co.jp/players/category/infielder"
	outfielderUrl := "https://www.yakult-swallows.co.jp/players/category/outfielder"

	pitchers := names(loadDocument(pitcherUrl))
	catchers := names(loadDocument(catcherUrl))
	ielders := names(loadDocument(infielderUrl))
	outfielder := names(loadDocument(outfielderUrl))

	players := append(append(append(pitchers, catchers...), ielders...), outfielder...)
	for _, name := range players {
		fmt.Println(name)
	}
}
