package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bcongdon/youtube-patreon-finder/lib"
)

func main() {
	file := os.Args[1]
	subscriptions, err := lib.FromFile(file)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for _, s := range subscriptions {
		fmt.Printf("%s\t%s\n", s.Channel.Name(), s.PatreonURL)
	}
}
