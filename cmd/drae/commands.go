package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/squat/drae/pkg/drae"
)

func api(port int, l *log.Logger) {
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		word := r.URL.Path[5:]
		err := encode(w, word)
		if err != nil {
			l.Errorf("failed to define word: %v", err)
			if _, ok := err.(drae.NotFoundError); ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	l.Infof("starting drae on port %d", port)
	l.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func define(word string, l *log.Logger) {
	if err := encode(os.Stdout, word); err != nil {
		l.Fatalf("failed to define word: %v", err)
	}
}

func encode(writer io.Writer, word string) error {
	enc := json.NewEncoder(writer)
	response, err := drae.Define(word)
	if err != nil {
		fmt.Println(word)
		return err
	}

	enc.Encode(response)
	return nil
}
