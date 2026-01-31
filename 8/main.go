package main

import (
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

var ntpServer = "time.google.com"

func main() {
	t, err := ntp.Time(ntpServer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting time from %v server: %v\n", ntpServer, err)
		os.Exit(1)
	}
	fmt.Println("Current time:", t)
}
