// Package printip contains an HTTP Cloud Function.
package printip

import (
	// "encoding/json"
	"fmt"
	// "html"
	// "io"
	// "log"
	"net/http"
)

// PrintIP out the IP address that is passed
func PrintIP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	// fmt.Println("Endpoint Hit: homePage")
	key, ok := r.URL.Query()["ip"]
	ip := "none"
	if ok {
		ip = key[0]
	}
	fmt.Fprintf(w, "IP address is: ")
	fmt.Fprintf(w, ip)
}
