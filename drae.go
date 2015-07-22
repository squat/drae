package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := flag.Int("port", 4000, "The port to run RAE on")
	flag.Parse()
	args := flag.Args()
	command := args[0]

	if command == "define" {
		respond(os.Stdout, args[1])
		return
	}

	if command == "api" {
		api(*port)
		return
	}

	fmt.Println("Unrecognized command:", command)
}

func api(port int) {
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		word := r.URL.Path[5:]
		respond(w, word)
	})

	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func respond(writer io.Writer, word string) *Entry {
	enc := json.NewEncoder(writer)
	response := Scrape(strings.ToLower(word))

	enc.Encode(response)
	return response
}
