package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Player = struct {
	Name     string `json:"name"`
	Hometown string `json:"hometown"`
}

func loadDocument(url string) *goquery.Document {
	res, err := http.Get(url)
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

func players(doc *goquery.Document) []Player {
	players := make([]Player, 0)
	doc.Find(".item-avatar").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		detailDoc := loadDocument("https://www.yakult-swallows.co.jp/" + href)

		name := strings.ReplaceAll(s.Find(".item-title").Text(), "　", " ")

		re := regexp.MustCompile(`県$|府$|都$`)
		hometownPath := "#top_ > div > div.sect > div > article > div.box-profile > div > div.md-6-5 > div > table > tbody > tr:nth-child(3) > td:nth-child(4)"
		hometown := re.ReplaceAllString(detailDoc.Find(hometownPath).Text(), "")

		players = append(players, Player{Name: name, Hometown: hometown})
	})

	return players
}

func main() {
	managerUrl := "https://www.yakult-swallows.co.jp/players/category/manager"
	pitcherUrl := "https://www.yakult-swallows.co.jp/players/category/pitcher"
	catcherUrl := "https://www.yakult-swallows.co.jp/players/category/catcher"
	infielderUrl := "https://www.yakult-swallows.co.jp/players/category/infielder"
	outfielderUrl := "https://www.yakult-swallows.co.jp/players/category/outfielder"

	managers := players(loadDocument(managerUrl))
	pitchers := players(loadDocument(pitcherUrl))
	catchers := players(loadDocument(catcherUrl))
	infielders := players(loadDocument(infielderUrl))
	outfielders := players(loadDocument(outfielderUrl))

	players := append(append(append(append(managers, pitchers...), catchers...), infielders...), outfielders...)

	file, err := os.Create(`SwallowsPlayers.json`)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bytes, _ := json.MarshalIndent(players, "", "  ")
	file.Write(([]byte)(string(bytes)))

	fmt.Println("SUCCESS")
}
