package gasbuddy

import (
	"context"
	"errors"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// Connect to an Influx Database
// return influxdb Client or errors
func ConnectToInfluxDB(dbToken string, dbURL string) (influxdb2.Client, error) {

	if dbToken == "" {
		return nil, errors.New("INFLUXDB_TOKEN must be set")
	}

	if dbURL == "" {
		return nil, errors.New("INFLUXDB_URL must be set")
	}

	client := influxdb2.NewClient(dbURL, dbToken)

	// validate client connection health
	_, err := client.Health(context.Background())

	return client, err
}

// WriteToInfluxDB writes station price data to InfluxDB
func WriteToInfluxDB(client influxdb2.Client, org string, bucket string, stationName string, price float32) error {
	// Get a blocking write API instance for synchronous writes
	writeAPI := client.WriteAPIBlocking(org, bucket)

	// Create a point with measurement "gas_prices", tag "station", and field "price"
	point := influxdb2.NewPointWithMeasurement("gas_price").
		AddTag("source", stationName).
		AddField("price", price).
		SetTime(time.Now())

	// Write the point synchronously
	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
		return fmt.Errorf("write point error: %w", err)
	}

	// Ensure all points are flushed
	if err := writeAPI.Flush(context.Background()); err != nil {
		return fmt.Errorf("flush error: %w", err)
	}

	return nil
}
