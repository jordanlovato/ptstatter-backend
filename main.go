package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", CalculatePtStats)
	http.ListenAndServe(":9054", nil)
}
