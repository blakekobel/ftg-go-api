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
	City           string `json:"city"`
	Country        string `json:"country"`
	CreatedDt      string `json:"created_dt"`
	ExpirationDate string `json:"exp_dt"`
}

// PrintIP out the IP address that is passed
func PrintIP(w http.ResponseWriter, r *http.Request) {
	// This chunk of code receives the query from the API Gateway
	// This checks for the following argument in the url /?ip={IP Address or DomainName}
	key, ok := r.URL.Query()["ip"]
	urlarg := "none"
	if ok {
		urlarg = key[0]
	} else {
		fmt.Fprint(w, "You need to call the api with /?ip={IP Address or DomainName}")
	}

	ipin := "none"
	domname := "none"

	//This sectionchecks if the api gave a domain or an IP address
	//It sets the ip and domain to a variable or fails if niether returns anything
	//This first line does a look up by the IP address, If its a domain it goes to else
	addrnames, addrerr := net.LookupAddr(urlarg)
	if addrerr == nil {
		domname = addrnames[0]
		ipin = urlarg
	} else {
		//This line looks up the Domain name to get first IP address affiliated
		ipaddr, iperr := net.LookupIP(urlarg)
		if iperr == nil {
			domname = urlarg
			ipin = ipaddr[0].String()
		} else {
			fmt.Fprint(w, addrerr)
			fmt.Fprint(w, iperr)
			fmt.Fprint(w, "Niether the IP address or the Domain returned a value")
		}
	}

	//Here is a nested If statement that returns information about the Domain/IP address
	//The first if statement uses a WhoIs module to get the raw data about a domain.
	//The second if statement checks if that raw data can be parsed using another module.
	if raw, rawErr := whois.Whois(domname); rawErr != nil {
		//Print Error Message if there was an error in the first IP lookup
		fmt.Fprint(w, rawErr)
	} else {
		// fmt.Println(raw)
		if result, err := whoisparser.Parse(raw); err != nil {
			// Print error if there was an issue parsing the raw WhoIs IP data
			fmt.Fprint(w, err)
		} else {
			//If no errors occur, we populate the struct created earlier with the data we want to provide
			e4 := IPAddr{
				IPAddress:      ipin,
				Domain:         result.Domain.Domain,
				City:           result.Registrar.City,
				Country:        result.Registrar.Country,
				CreatedDt:      result.Domain.CreatedDate,
				ExpirationDate: result.Domain.ExpirationDate,
			}
			//Write that JSON back to the URL
			json.NewEncoder(w).Encode(e4)
		}
	}
}
