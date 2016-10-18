package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

var (
	procdeckdata map[string]string
	procstatdata map[string]map[string]*Result
	decklists    map[string]string
)

func CalculatePtStats(w http.ResponseWriter, r *http.Request) {
	staturls := []string{
		"http://magic.wizards.com/en/events/coverage/ptkld/round-4-results-2016-10-14",
		"http://magic.wizards.com/en/events/coverage/ptkld/round-5-results-2016-10-14",
		"http://magic.wizards.com/en/events/coverage/ptkld/round-6-results-2016-10-14",
		"http://magic.wizards.com/en/events/coverage/ptkld/round-7-results-2016-10-14",
		"http://magic.wizards.com/en/events/coverage/ptkld/round-8-results-2016-10-14",
		"http://magic.wizards.com/en/events/coverage/ptkld/round-12-results-2016-10-14",
		"http://magic.wizards.com/en/events/coverage/ptkld/round-13-results-2016-10-14",
		"http://magic.wizards.com/en/events/coverage/ptkld/round-14-results-2016-10-14",
		"http://magic.wizards.com/en/events/coverage/ptkld/round-15-results-2016-10-14",
		"http://magic.wizards.com/en/events/coverage/ptkld/round-16-results-2016-10-14",
	}

	deckurls := []string{
		"http://magic.wizards.com/en/events/coverage/ptkld/24-27-point-standard-decklists-2016-10-16",
		"http://magic.wizards.com/en/events/coverage/ptkld/21-23-point-standard-decklists-2016-10-16",
		"http://magic.wizards.com/en/events/coverage/ptkld/18-20-point-standard-decklists-2016-10-16",
	}

	type Counter struct {
		Ch    chan int
		Count int
	}

	statcnt := &Counter{
		Ch:    make(chan int, len(staturls)),
		Count: 0,
	}

	statdata := make([]string, 0)
	for _, u := range staturls {
		go func(url string, c *Counter) {
			arr := ScrapeStats(url)
			statdata = append(statdata, arr...)
			c.Count += 1
			c.Ch <- c.Count
			if c.Count == len(staturls) {
				close(c.Ch)
			}
		}(u, statcnt)
	}

	deckcnt := &Counter{
		Ch:    make(chan int, len(deckurls)),
		Count: 0,
	}

	deckdata := make([]string, 0)
	for _, u := range deckurls {
		go func(url string, c *Counter) {
			arr := ScrapeDecks(url)
			deckdata = append(deckdata, arr...)
			c.Count += 1
			c.Ch <- c.Count
			if c.Count == len(deckurls) {
				close(c.Ch)
			}
		}(u, deckcnt)
	}

	ready := make(chan int)
	for i := range statcnt.Ch {
		if i == len(staturls) {
			go BuildResultLookup(statdata, ready)
		}
	}

	for i := range deckcnt.Ch {
		if i == len(deckurls) {
			go BuildDecksLookup(deckdata, ready)
		}
	}

	for i := range ready {
		if i == 1 {
			// Stats calculated, encode to json and ship it back
			var Data struct {
				D map[string]map[string]*Result `json:"data"`
				L map[string]string             `json:"legend"`
			}
			Data.D = procstatdata
			Data.L = decklists

			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(Data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Accept, Accept-Language, Content-Language, Content-Type")
			w.WriteHeader(http.StatusOK)

			if _, err := io.Copy(w, &buf); err != nil {
				panic(err.Error())
			}
		}
	}
}
