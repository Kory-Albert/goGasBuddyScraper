package main

import (
	"encoding/json"
	"fmt"
	"gasbuddyscraper/gasbuddy"
	"log"
	"os"
	"strings"
)

// Environment Variables
var (
	influxdbToken = os.Getenv("INFLUXDB_TOKEN")
	influxURL     = os.Getenv("INFLUXDB_URL")
	influxOrg     = os.Getenv("INFLUXDB_ORG")
	influxBucket  = os.Getenv("INFLUXDB_BUCKET")

	stationNames = os.Getenv("STATION_NAMES")
	stations     = os.Getenv("STATION_IDS")
	stringIDs    = strings.Split(stations, ",")
	stringNames  = strings.Split(stationNames, ",")
)

func main() {
	// Check 1:1 match on Station ID/Names
	if len(stringIDs) != len(stringNames) {
		log.Fatal("STATION_IDS and STATION_NAMES must be 1:1 match")
	}
	client := gasbuddy.NewClient()

	fmt.Println("Connecting to InfluxDB...")
	conn, err := gasbuddy.ConnectToInfluxDB(influxdbToken, influxURL)
	if err != nil {
		fmt.Println("Error connecting to InfluxDB!")
		log.Fatal(err)
	}
	fmt.Println("Connected to InfluxDB successfully")
	defer conn.Close()

	for i := range stringIDs {
		raw, err := client.GetStationPrices(stringIDs[i])
		if err != nil {
			log.Fatal(err)
		}

		var parsed gasbuddy.StationResponse
		if err := json.Unmarshal(raw, &parsed); err != nil {
			log.Fatal(err)
		}

		for _, p := range parsed.Data.Station.Prices {
			if p.LongName == "Regular" {
				fmt.Printf("Extracted price for %s: %s\n", stringNames[i], p.Credit.FormattedPrice)
				fmt.Printf("Writing to InfluxDB: station=%s, price=%.2f\n", stringNames[i], p.Credit.Price)
				// Write to InfluxDB
				err := gasbuddy.WriteToInfluxDB(conn, influxOrg, influxBucket, stringNames[i], p.Credit.Price)
				if err != nil {
					log.Printf("Error writing to InfluxDB: %v", err)
				} else {
					fmt.Println("Successfully wrote to InfluxDB")
				}
			}
		}
	}
}
