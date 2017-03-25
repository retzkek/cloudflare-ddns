// Dynamically update CloudFlare zone file with current IP address
//
// Copyright 2013 Kevin Retzke <kmr@kmr.me>
// Shared under MIT license, see LICENSE for details
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	APIURL    = "https://www.cloudflare.com/api_json.html"
	EMAIL     = "nobody@example.com"
	TKN       = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	ZONE      = "example.com"
	DOMAIN    = "example.com"
	EGRESS    = "re0"
	FREQUENCY = 5 * time.Minute
)

// types for unmarshalling the API rec_load_all response

type ApiRecObj struct {
	Rec_id    string
	Zone_name string
	Name      string
	Type      string
	Content   string
	// don't need the rest of the fields
}

type ApiRecs struct {
	Has_more bool
	Count    int
	Objs     []ApiRecObj
}

type ApiResponse struct {
	Recs ApiRecs
}

type ApiRecLoadAll struct {
	Request  map[string]string
	Response ApiResponse
	Result   string
	Msg      string
}

// Get the domain's A record
func getRec() (*ApiRecObj, error) {
	args := url.Values{}
	args.Set("a", "rec_load_all")
	args.Set("tkn", TKN)
	args.Set("email", EMAIL)
	args.Set("z", ZONE)
	resp, err := http.PostForm(APIURL, args)
	if err != nil {
		return nil, fmt.Errorf("Error posting request: %v", err)
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	var m ApiRecLoadAll
	err = dec.Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("Error decoding response: %v", err)
	}
	if m.Result != "success" {
		return nil, fmt.Errorf("API call returned error: %v", m.Msg)
	}
	for _, rec := range m.Response.Recs.Objs {
		if rec.Type == "A" && rec.Name == DOMAIN {
			return &rec, nil
		}
	}
	return nil, fmt.Errorf("Unable to locate rec ID")
}

// Get the current IP address from onhub api.
func getAddr(ifname string) string {
	r, err := http.Get("http://onhub.here/api/v1/status")
	if err != nil {
		log.Printf("error contacting onhub: %s\n", err)
		return ""
	}
	var stat OnhubStatus
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&stat)
	if err != nil {
		log.Printf("error unmarshalling onhub response: %s\n", err)
		return ""
	}
	if !stat.Wan.Online {
		log.Printf("error: onhub not online")
		return ""
	}
	return stat.Wan.LocalIpAddress
}

// Get current address on record, update if different.
func updateAddr(newaddr string) error {
	rec, err := getRec()
	if err != nil {
		return fmt.Errorf("Error getting DNS record: %v", err)
	}
	if rec.Content == newaddr {
		log.Printf("DNS record matches current IP\n")
		return nil
	}
	args := url.Values{}
	args.Set("a", "rec_edit")
	args.Set("tkn", TKN)
	args.Set("email", EMAIL)
	args.Set("z", ZONE)
	args.Set("type", "A")
	args.Set("id", rec.Rec_id)
	args.Set("name", rec.Name)
	args.Set("content", newaddr)
	args.Set("ttl", fmt.Sprintf("%.0f", FREQUENCY.Seconds())) // 1=Automatic, otherwise set between 120 and 4,294,967,295 seconds
	args.Set("service_mode", "0")                             // 1 = orange cloud, 0 = grey cloud
	resp, err := http.PostForm(APIURL, args)
	if err != nil {
		return fmt.Errorf("Error posting request: %v", err)
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	// not exactly right, but ApiRecLoadAll will get us the result and msg
	var m ApiRecLoadAll
	err = dec.Decode(&m)
	if err != nil {
		return fmt.Errorf("Error decoding response: %v", err)
	}
	if m.Result != "success" {
		return fmt.Errorf("API call returned error: %v", m.Msg)
	}
	log.Printf("Successfully updated DNS record.\n")
	return nil
}

func main() {
	log.Printf("Starting DDNS updates for %v @ frequency %v\n", DOMAIN, FREQUENCY)
	lastip := getAddr(EGRESS)
	if lastip == "" {
		// initial errors are fatal since we want to ensure a clean starting state
		log.Fatalf("Unable to determine current IP address.\n")
	} else {
		log.Printf("Current IP: %v\n", lastip)
	}
	if err := updateAddr(lastip); err != nil {
		log.Fatal(err)
	}
	ch := time.Tick(FREQUENCY)
	for _ = range ch {
		if a := getAddr(EGRESS); a != lastip {
			if a == "" {
				// errors within loop are not fatal, we'll just try again next go round
				log.Printf("Unable to determine IP address.\n")
			} else {
				log.Printf("IP address changed. New IP: %v\n", a)
				if err := updateAddr(a); err != nil {
					log.Printf("Error updating DNS record: %v\n", err)
				} else {
					lastip = a
				}
			}
		}
	}
}
