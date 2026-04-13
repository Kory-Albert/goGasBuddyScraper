# GoGasBuddyScraper

A Go application that scrapes gas prices from GasBuddy and stores them in InfluxDB.

## Overview

This project scrapes regular gas prices from specified GasBuddy stations and writes the data to an InfluxDB database for monitoring and analysis. This is a replacement for my older Python scraper that used Beautiful Soup. This Go version will create the necessary cookies and headers to use the GraphQL API for quicker results.

## Features

- Scrapes gas prices from GasBuddy using their GraphQL API
- Extracts regular gas prices for specified stations
- Writes data to InfluxDB with proper timestamping
- Environment variable configuration for flexibility
- Error handling and logging

## Environment Variables

The following environment variables must be set before running the application:

> STATION_IDS and STATION_NAMES must be mapped 1:1

| Variable | Description | Example |
|----------|-------------|---------|
| `INFLUXDB_TOKEN` | InfluxDB authentication token | `aK9nY3bZ6x8T0W5YhJ4qM1sO3vC7gR9pZbF6dN8uLwQ3oS2mH4iT5uJ9kG1vXq0lY==` |
| `INFLUXDB_URL` | InfluxDB server URL | `https://influxdb.local` |
| `INFLUXDB_ORG` | InfluxDB organization name | `home` |
| `INFLUXDB_BUCKET` | InfluxDB bucket name | `gas_prices` |
| `STATION_IDS` | Comma-separated list of GasBuddy station IDs | `177667,19606` |
| `STATION_NAMES` | Comma-separated list of station names (must match order of STATION_IDS) | `Roseville-Costco,Another-Station` |

## Installation

1. Clone the repository
2. Install Go dependencies:
   ```bash
   go mod download
   ```
3. Set the required environment variables
4. Run the application:
   ```bash
   go run main.go
   ```

## How It Works

1. The application connects to InfluxDB using the provided credentials
2. For each station ID in `STATION_IDS`:
   - Fetches the CSRF token from GasBuddy
   - Queries the GasBuddy GraphQL API for station prices
   - Extracts the regular gas price
   - Writes the price to InfluxDB with:
     - Measurement: `gas_prices`
     - Tag: `source` (station name)
     - Field: `price` (numeric price value)
     - Timestamp: Current time

## Data Structure in InfluxDB

Each data point written to InfluxDB has:
- **Measurement**: `gas_prices`
- **Tags**:
  - `source`: The name of the gas station (from STATION_NAMES)
- **Fields**:
  - `price`: The gas price as a float64
- **Timestamp**: When the data was collected

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This project is for educational purposes only. Please check GasBuddy's terms of service before using this scraper in production.
