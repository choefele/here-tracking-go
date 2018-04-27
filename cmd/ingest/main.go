package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/choefele/here-tracking-go/pkg/tracking"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ingest device_id device_secret")
		os.Exit(-1)
	}

	client := tracking.NewDeviceClient(os.Args[1], os.Args[2])
	dr := &tracking.DataRequest{
		Timestamp: tracking.Time{Time: time.Now()},
		Position: &tracking.Position{
			Lat:      52,
			Lng:      13,
			Accuracy: 100,
		},
	}
	err := client.Ingestion.Send(context.Background(), []*tracking.DataRequest{dr})
	fmt.Printf("Send: done, error: %v\n", err)
}
