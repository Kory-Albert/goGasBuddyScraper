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
var influxdbToken string = os.Getenv("INFLUXDB_TOKEN")
var influxURL string = os.Getenv("INFLUXDB_URL")
var influxOrg string = os.Getenv("INFLUXDB_ORG")
var influxBucket string = os.Getenv("INFLUXDB_BUCKET")

var stationNames string = os.Getenv("STATION_NAMES")
var stations string = os.Getenv("STATION_IDS")

func main() {
	client := gasbuddy.NewClient()

	fmt.Println("Connecting to InfluxDB...")
	conn, err := gasbuddy.ConnectToInfluxDB(influxdbToken, influxURL)
	if err != nil {
		fmt.Println("Error connecting to InfluxDB!")
		log.Fatal(err)
	}
	fmt.Println("Connected to InfluxDB successfully")
	defer conn.Close()

	var stringIDs []string = strings.Split(stations, ",")
	var stringNames []string = strings.Split(stationNames, ",")

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
