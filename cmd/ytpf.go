package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/bcongdon/youtube-patreon-finder/pkg/ytpf"
	"github.com/olekukonko/tablewriter"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("USAGE: ytpf <opml_file>")
		os.Exit(1)
	}
	file := os.Args[1]
	subscriptions, err := ytpf.FromOPMLFile(file)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Sort subscriptions alphabetically by name
	sort.Slice(subscriptions, func(a, b int) bool {
		x := strings.ToLower(subscriptions[a].Channel.Name())
		y := strings.ToLower(subscriptions[b].Channel.Name())
		return x < y
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Channel", "Patreon URL"})

	for _, s := range subscriptions {
		if s.PatreonURL == "" {
			continue
		}
		table.Append([]string{s.Channel.Name(), s.PatreonURL})
	}
	table.Render()
}
