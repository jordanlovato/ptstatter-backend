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
		if n.DataAtom == atom.Span && n.Parent != nil && n.Parent.Parent != nil {
			return scrape.Attr(n, "class") == "deck-meta"
		}
		return false
	}

	decks := scrape.FindAll(doc, matcher)
	for _, deck := range decks {
		p := deck.Parent
		h4, _ := scrape.Find(deck, func(n *html.Node) bool {
			if n.DataAtom == atom.H4 {
				return true
			}
			return false
		})

		split := strings.SplitN(scrape.Text(h4), "-", 2)
		uri := scrape.Attr(p, "id")
		deckurl := sourceUrl + "#" + uri
		decklookup = append(decklookup, strings.TrimSpace(split[0]), strings.TrimSpace(split[1]), deckurl)
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
			statlookup = append(statlookup, scrape.Text(stat))
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
