package main

import (
	"context"
	"fmt"
	"os"

	"github.com/choefele/here-tracking-go/pkg/tracking"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: admin email password")
		os.Exit(-1)
	}

	client := tracking.NewAdminClient(os.Args[1], os.Args[2])
	err := client.User.ListDevices(context.Background())
	fmt.Printf("ListDevices: done, error: %v\n", err)
}
