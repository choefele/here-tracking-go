package main

import (
	"context"
	"fmt"
	"os"

	"github.com/choefele/here-tracking-go/pkg/tracking"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: user email password")
		os.Exit(-1)
	}

	client := tracking.NewClient("", "")

	t, err := client.User.Login(context.Background(), os.Args[1], os.Args[2])
	fmt.Printf("Login: %v, error: %v\n", t, err)
	if err != nil {
		os.Exit(-1)
	}

	client.AccessToken = &t.AccessToken // fake
	client.UserAccessToken = &t.AccessToken
	err = client.User.ListDevices(context.Background())
	fmt.Printf("ListDevices: done, error: %v\n", err)
}
