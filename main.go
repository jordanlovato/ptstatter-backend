package main

import (
	"fmt"
	"math"
)

func main() {
	ss := make([][]string, 10)
	ss[0] = ScrapeStats("http://magic.wizards.com/en/events/coverage/ptkld/round-4-results-2016-10-14")
	ss[1] = ScrapeStats("http://magic.wizards.com/en/events/coverage/ptkld/round-5-results-2016-10-14")
	ss[2] = ScrapeStats("http://magic.wizards.com/en/events/coverage/ptkld/round-6-results-2016-10-14")
	ss[3] = ScrapeStats("http://magic.wizards.com/en/events/coverage/ptkld/round-7-results-2016-10-14")
	ss[4] = ScrapeStats("http://magic.wizards.com/en/events/coverage/ptkld/round-8-results-2016-10-14")
	ss[5] = ScrapeStats("http://magic.wizards.com/en/events/coverage/ptkld/round-12-results-2016-10-14")
	ss[6] = ScrapeStats("http://magic.wizards.com/en/events/coverage/ptkld/round-13-results-2016-10-14")
	ss[7] = ScrapeStats("http://magic.wizards.com/en/events/coverage/ptkld/round-14-results-2016-10-14")
	ss[8] = ScrapeStats("http://magic.wizards.com/en/events/coverage/ptkld/round-15-results-2016-10-14")
	ss[9] = ScrapeStats("http://magic.wizards.com/en/events/coverage/ptkld/round-16-results-2016-10-14")
	var stats []string
	for _, s := range ss {
		stats = append(stats, s...)
	}
	d1 := ScrapeDecks("http://magic.wizards.com/en/events/coverage/ptkld/24-27-point-standard-decklists-2016-10-16")
	d2 := ScrapeDecks("http://magic.wizards.com/en/events/coverage/ptkld/21-23-point-standard-decklists-2016-10-16")
	d3 := ScrapeDecks("http://magic.wizards.com/en/events/coverage/ptkld/18-20-point-standard-decklists-2016-10-16")
	d := append(d1, d2...)
	d = append(d, d3...)
	lookup := BuildDecksLookup(d)
	data := BuildResultLookup(stats, lookup)
	fmt.Printf("%#v\n", data)
}

type Result struct {
	WinCnt int
	LossCnt int
	Total int
	WinRate float64
}

func BuildResultLookup(stats []string, decks map[string]string) map[string]map[string]*Result {
	winloss := map[string]map[string]*Result{}
	for i, _ := range stats {
		j := i % 3
		if j == 0 {
			player := stats[j]
			playerdeck := decks[player]
			result := stats[j+1]
			var resultstr, matchscore string
			fmt.Sscanf(result, "%s %s", &resultstr, &matchscore)
			fmt.Println("RESULT: ", resultstr, matchscore)
			opponent := stats[j+2]
			opponentdeck := decks[opponent]

			var defaultwincnt, defaultlosscnt int
			if resultstr == "won" {
				defaultwincnt = 1
				defaultlosscnt = 0
			} else if resultstr == "lost" {
				defaultwincnt = 0
				defaultlosscnt = 1
			}

			if _, ok := winloss[playerdeck]; !ok {
				// if the player does not exist in the lookup
				winloss[playerdeck] = map[string]*Result{}
				winloss[playerdeck][opponentdeck] = &Result{WinCnt:defaultwincnt, LossCnt:defaultlosscnt, Total:1}
			} else {
				if _, ok := winloss[playerdeck][opponentdeck]; !ok {
					// if the player exists in the lookup, but not the opposing deck
					winloss[playerdeck][opponentdeck] = &Result{WinCnt:defaultwincnt, LossCnt:defaultlosscnt, Total:1}
				} else {
					// if all keys exist
					if resultstr == "won" {
						winloss[playerdeck][opponentdeck].WinCnt++
						winloss[playerdeck][opponentdeck].Total++
					} else if resultstr == "lost" {
						winloss[playerdeck][opponentdeck].LossCnt++
						winloss[playerdeck][opponentdeck].Total++
					}
				}
			}
		}
	}

	for _, matchups := range winloss {
		for _, res := range matchups {
			rate := float64(res.WinCnt) / float64(res.Total)
			res.WinRate = math.Floor(rate)
		}
	}
	fmt.Printf("%#v\n", winloss)

	return winloss

}

func BuildDecksLookup(decks []string) map[string]string {
	deckslookup := map[string]string{}
	for i, deck := range decks {
		j := i % 3
		switch (j) {
		case 0:
			deckslookup[deck] = decks[j+1]
		}
	}
	return deckslookup
}
