// Dynamically update CloudFlare zone file with current IP address
//
// Copyright 2013-2018 Kevin Retzke <kmr@kmr.me>
// Shared under MIT license, see LICENSE for details
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	cloudflare "github.com/cloudflare/cloudflare-go"
)

const (
	EMAIL     = "nobody@example.com"
	TKN       = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	ZONE      = "example.com"
	DOMAIN    = "example.com"
	FREQUENCY = 5 * time.Minute
)

// Get the current IP address from onhub api.
func getAddr() string {
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
func updateAddr(api *cloudflare.API, zoneID string, newaddr string) error {
	filter := cloudflare.DNSRecord{
		Name: DOMAIN,
		Type: "A",
	}
	recs, err := api.DNSRecords(zoneID, filter)
	if err != nil {
		return err
	}
	if len(recs) != 1 {
		return fmt.Errorf("got %d A records for %s (expected 1)", len(recs), DOMAIN)
	}
	oldaddr := recs[0].Content
	if oldaddr == newaddr {
		log.Printf("current address %s same as new %s", oldaddr, newaddr)
		return nil
	}
	if err = api.UpdateDNSRecord(zoneID, recs[0].ID, cloudflare.DNSRecord{Content: newaddr}); err != nil {
		return fmt.Errorf("error updating DNS record: %s", err.Error())
	}
	log.Printf("Successfully updated DNS record.\n")
	return nil
}

func main() {
	log.Printf("Starting DDNS updates for %v @ frequency %v\n", DOMAIN, FREQUENCY)
	lastip := getAddr()
	if lastip == "" {
		// initial errors are fatal since we want to ensure a clean starting state
		log.Fatalf("Unable to determine current IP address.\n")
	} else {
		log.Printf("Current IP: %v\n", lastip)
	}

	// load cloudflare API library
	api, err := cloudflare.New(TKN, EMAIL)
	if err != nil {
		log.Fatalf("error loading cloudflare APi: %s\n", err.Error())
	}

	// Fetch the zone ID
	id, err := api.ZoneIDByName(ZONE) // Assuming example.com exists in your Cloudflare account already
	if err != nil {
		log.Fatalf("error getting zone ID: %s\n", err.Error())
	}
	log.Printf("zone ID: %s\n", id)

	if err := updateAddr(api, id, lastip); err != nil {
		log.Fatal(err)
	}
	ch := time.Tick(FREQUENCY)
	for _ = range ch {
		a := getAddr()
		if a != lastip {
			if a == "" {
				// errors within loop are not fatal, we'll just try again next go round
				log.Printf("Unable to determine IP address.\n")
			} else {
				log.Printf("IP address changed. New IP: %v\n", a)
				if err := updateAddr(api, id, a); err != nil {
					log.Printf("Error updating DNS record: %v\n", err)
				} else {
					lastip = a
				}
			}
		} else {
			log.Print("address unchanged")
		}
	}
}
