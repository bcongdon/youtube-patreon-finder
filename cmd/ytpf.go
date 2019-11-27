package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/bcongdon/youtube-patreon-finder/pkg/ytpf"
	"github.com/olekukonko/tablewriter"
)

var serverFlag = flag.Bool("server", false, "Start server")
var portFlag = flag.String("port", "8080", "Port to listen on")

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] \n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
}

func cli() {
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
		x := strings.ToLower(subscriptions[a].Channel.Name)
		y := strings.ToLower(subscriptions[b].Channel.Name)
		return x < y
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Channel", "Patreon URL"})

	for _, s := range subscriptions {
		if s.PatreonURL == "" {
			continue
		}
		table.Append([]string{s.Channel.Name, s.PatreonURL})
	}
	table.Render()
}

func server() {
	handler := &ytpf.Handler{}
	log.Fatal(http.ListenAndServe(":"+*portFlag, handler))
}

func main() {
	if *serverFlag {
		server()
	} else {
		cli()
	}
}
