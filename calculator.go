package main

import (
	"fmt"
	"math"
)

type Result struct {
	WinCnt  int     `json:"wins"`
	LossCnt int     `json:"losses"`
	Total   int     `json:"total_matches"`
	WinRate float64 `json:"win_rate"`
}

func BuildResultLookup(stats []string, ch chan int) {
	// sync with decks lookup
	<-ch
	// do calc
	winloss := map[string]map[string]*Result{}
	for i, _ := range stats {
		j := i % 3
		if j == 0 {
			k := i + j
			player := stats[k]
			playerdeck := procdeckdata[player]
			result := stats[k+1]
			var resultstr, matchscore string
			fmt.Sscanf(result, "%s %s", &resultstr, &matchscore)
			opponent := stats[k+2]
			opponentdeck := procdeckdata[opponent]
			if player == "" || playerdeck == "" || result == "" || opponent == "" || opponentdeck == "" {
				//fmt.Printf("[BAD DATA] \nPLAYER: %s\nPLAYERDECK: %s\nRESULT: %s\nOPPONENT: %s\n OPPONENTDECK: %s\n", player, playerdeck, result, opponent, opponentdeck)
				continue
			}

			// initialized any uninitialized decks
			matchupdecks := []string{playerdeck, opponentdeck}
			for i, initialdeck := range matchupdecks {
				conversedeck := matchupdecks[(i+1)%2]
				if _, ok := winloss[initialdeck]; !ok {
					// if the player does not exist in the lookup
					winloss[initialdeck] = map[string]*Result{}
					winloss[initialdeck][conversedeck] = &Result{WinCnt: 0, LossCnt: 0, Total: 0}
				} else {
					if _, ok := winloss[initialdeck][conversedeck]; !ok {
						// if the player exists in the lookup, but not the opposing deck
						winloss[initialdeck][conversedeck] = &Result{WinCnt: 0, LossCnt: 0, Total: 0}
					}
				}
			}

			if resultstr == "won" {
				winloss[playerdeck][opponentdeck].WinCnt++
				winloss[playerdeck][opponentdeck].Total++
				winloss[opponentdeck][playerdeck].LossCnt++
				winloss[opponentdeck][playerdeck].Total++
			} else if resultstr == "lost" {
				winloss[playerdeck][opponentdeck].LossCnt++
				winloss[playerdeck][opponentdeck].Total++
				winloss[opponentdeck][playerdeck].WinCnt++
				winloss[opponentdeck][playerdeck].Total++
			}
		}
	}

	for _, matchups := range winloss {
		for _, res := range matchups {
			rate := (float64(res.WinCnt) / float64(res.Total)) * 100
			if math.IsNaN(rate) {
				res.WinRate = float64(0)
			} else {
				res.WinRate = math.Floor(rate)
			}
		}
	}

	procstatdata = winloss
	ch <- 1
	close(ch)
}

func BuildDecksLookup(decks []string, ch chan int) {
	procdeckdata = make(map[string]string)
	decklists = make(map[string]string)
	for i, deck := range decks {
		j := i % 3
		switch j {
		case 0:
			k := i + j
			procdeckdata[deck] = decks[k+1]
			if _, exists := decklists[decks[k+1]]; !exists {
				decklists[decks[k+1]] = ""
			}
		}
	}
	ch <- 0
}
