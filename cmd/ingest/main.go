package main

import (
	"context"
	"fmt"

	"github.com/choefele/here-tracking-go/pkg/tracking"
)

func main() {
	client := tracking.NewClient()
	h, e := client.Ingestion.Health(context.Background())
	fmt.Printf("Health: %v, error: %v\n", h, e)

	dr := &tracking.DataRequest{
		Position: &tracking.Position{
			Lat:      52,
			Lng:      13,
			Accuracy: 100,
		},
	}
	d, e := client.Ingestion.Send(context.Background(), dr)
	fmt.Printf("Send: %v, error: %v\n", d, e)
}
