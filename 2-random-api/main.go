package main

import (
	"math/rand"
	"net/http"
	"strconv"
)

func randRange(min, max int) string {
	return strconv.Itoa(rand.Intn(max-min+1) + min)
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	randNum := randRange(1, 6)
	w.Write([]byte(randNum))
	return
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.ListenAndServe(":8081", nil)
}
