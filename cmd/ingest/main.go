package main

import (
	"context"
	"fmt"
	"time"

	"github.com/choefele/here-tracking-go/pkg/tracking"
)

func main() {
	client := tracking.NewClient(
		"9d9c31be-dd5d-40b1-95af-7d5375c39561",
		"vHrFUhnxo0hxw2VqR5OXBBnvjeTK0T8etmws8HZ9dvw",
	)
	h, e := client.Ingestion.Health(context.Background())
	fmt.Printf("Health: %v, error: %v\n", h, e)

	dr := &tracking.DataRequest{
		Timestamp: tracking.Time{Time: time.Now()},
		Position: &tracking.Position{
			Lat:      52,
			Lng:      13,
			Accuracy: 100,
		},
	}
	e = client.Ingestion.Send(context.Background(), []*tracking.DataRequest{dr})
	fmt.Printf("Send: done, error: %v\n", e)
}
