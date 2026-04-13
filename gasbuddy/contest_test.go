package gasbuddy

import (
	"context"
	"testing"
)

// Test_connectToInfluxDB tests the connection to InfluxDB using hardcoded credentials
func Test_connectToInfluxDB(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Successful connection to InfluxDB",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConnectToInfluxDB("token", "https://influxdb.url")
			if (err != nil) != tt.wantErr {
				t.Errorf("ConnectToInfluxDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				health, healthErr := got.Health(context.Background())
				if healthErr != nil || health.Status != "pass" {
					t.Errorf("connectToInfluxDB() health check failed: err=%v, status=%v", healthErr, health.Status)
				}
				got.Close()
			}
		})
	}
}
