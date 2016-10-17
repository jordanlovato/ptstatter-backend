package main

import (
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"net/http"
	"strings"
)

func ScrapeDecks(sourceUrl string) []string {
	decklookup := make([]string, 0)
	b := makeRequestAndReturnBody(sourceUrl)
	doc, _ := html.Parse(b)
	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.Div && n.Parent != nil && n.Parent.Parent != nil {
			return scrape.Attr(n, "class") == "deck-group"
		}
		return false
	}

	decks := scrape.FindAll(doc, matcher)
	for _, deck := range decks {
		uri := scrape.Attr(deck, "id")
		split := strings.Split(uri, "_-_")
		deckurl := sourceUrl + "#" + uri
		decklookup = append(decklookup, strings.Replace(strings.TrimSpace(strings.ToLower(split[0])), "_", ", ", 1), strings.TrimSpace(strings.ToLower(split[1])), deckurl)
	}
	return decklookup
}

func ScrapeStats(sourceUrl string) []string {
	statlookup := make([]string, 0)
	b := makeRequestAndReturnBody(sourceUrl)
	doc, _ := html.Parse(b)
	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.Td && n.Parent != nil && n.Parent.Parent != nil && n.Parent.Parent.Parent != nil {
			return scrape.Attr(n.Parent.Parent.Parent, "class") == "sortable-table"
		}
		return false
	}
	stats := scrape.FindAll(doc, matcher)
	for i, stat := range stats {
		j := i % 7
		switch j {
		case 1, 3, 5:
			statlookup = append(statlookup, strings.ToLower(scrape.Text(stat)))
		}
	}
	return statlookup
}

func makeRequestAndReturnBody(url string) io.Reader {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := client.Do(req)
	return resp.Body
}
