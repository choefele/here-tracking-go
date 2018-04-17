package main

import (
	"context"
	"fmt"

	"github.com/choefele/here-tracking-go/pkg/tracking"
)

func main() {
	client := tracking.NewClient()
	h, e := client.Ingestion.Health(context.Background())
	fmt.Printf("Health: %v, error: %v\n", h.Message, e)
}
