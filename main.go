package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Wait(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()["w"]
	wait := "0" //no wait by default

	if len(q) > 0 {
		wait = q[0]
	}

	s, err := strconv.Atoi(wait)
	if err != nil {
		fmt.Fprint(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

	time.Sleep(time.Duration(s) * time.Second)
	fmt.Fprintf(w, fmt.Sprintf("I slept %d seconds!", s))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", Wait)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
