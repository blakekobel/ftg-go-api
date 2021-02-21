// Package printip contains the github modules used below.
package printip

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/likexian/whois-go"
	whoisparser "github.com/likexian/whois-parser-go"
)

// IPAddr ...
type IPAddr struct {
	IPAddress      string `json:"IP"`
	Domain         string `json:"Domain"`
	CreatedDt      string `json:"created_dt"`
	ExpirationDate string `json:"exp_dt"`
}

// PrintIP out the IP address that is passed
func PrintIP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	// fmt.Println("Endpoint Hit: homePage")
	key, ok := r.URL.Query()["ip"]
	ipin := "none"
	if ok {
		ipin = key[0]
	}

	names, err := net.LookupAddr(ipin)
	if err == nil {
		fmt.Println(names[0])
	} else {
		fmt.Println(err)
	}
	//This if statement queries the WHO
	if raw, rawErr := whois.Whois(names[0]); rawErr != nil {
		//Print Error Message if rawErr is not null
		fmt.Fprint(w, rawErr)
	} else {
		//
		// fmt.Println(raw)
		if result, err := whoisparser.Parse(raw); err != nil {
			fmt.Fprint(w, err)
		} else {
			// fmt.Println(result.Domain.Domain)
			// fmt.Println(result.Domain.ExpirationDate)

			// // Print the registrar name
			// fmt.Println(result.Technical.City)

			// // Print the registrant name
			// fmt.Println(result.Registrant.Name)

			// Print the registrant email address
			// fmt.Println(result.Registrant.Email)
			e4 := IPAddr{
				IPAddress:      ipin,
				Domain:         result.Domain.Domain,
				CreatedDt:      result.Domain.CreatedDate,
				ExpirationDate: result.Domain.ExpirationDate,
			}
			json.NewEncoder(w).Encode(e4)
		}
	}
}
