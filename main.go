package main

import (
	"fmt"
	"log"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
	key, ok := r.URL.Query()["ip"]
	ip := "none"
	if ok {
		ip = key[0]
	}
	fmt.Fprintf(w, "IP address is: ")
	fmt.Fprintf(w, ip)
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	fmt.Println("Server Started at Port 8000")
	// http.HandleFunc("/216.239.32.21/articles", returnAllArticles)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func main() {
	handleRequests()
}
