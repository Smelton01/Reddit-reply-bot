package scrape

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/PuerkitoBio/goquery"
)

const data = "./streamer-data.html"

type streamer struct {
	Rank  int
	Name  string
	Url   string
	Money string
}

func GetData() []streamer {
	streamerLeaderboard := []streamer{}
	html, err := ioutil.ReadFile(data)
	if err != nil {
		log.Fatalf("failed to read file %v: %v", data, err)
	}
	r := bytes.NewReader(html)

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatalf("failed to open document: %v", err)
	}
	doc.Find("div.Profile_Card").Each(func(i int, s *goquery.Selection) {
		strimer := streamer{}
		strimer.Name = s.Find("a").Text()
		strimer.Rank = i + 1
		strimer.Money = s.Children().Last().Text()
		strimer.Url, _ = s.Find("a").Attr("href")

		streamerLeaderboard = append(streamerLeaderboard, strimer)
	})

	return streamerLeaderboard

}
