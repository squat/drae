package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/squat/drae/pkg/version"
)

func main() {
	flags := struct {
		logLevel string
		port     int
		version  bool
	}{}

	flag.StringVarP(&flags.logLevel, "log-level", "l", "warn", "log level, e.g. \"debug\"")
	flag.IntVarP(&flags.port, "port", "p", 4000, "the port on which to run drae")
	flag.BoolVarP(&flags.version, "version", "v", false, "print version and exit")
	flag.Parse()

	if len(os.Args) == 1 {
		usage()
	}

	lvl, err := log.ParseLevel(flags.logLevel)
	if err != nil {
		log.Fatalf("invalid log-level: %v", err)
	}

	command := os.Args[1]
	if flags.version || command == "version" {
		fmt.Println(version.Version)
		return
	}

	l := log.New()
	l.Level = lvl

	switch command {
	case "define":
		if len(os.Args) == 2 {
			usageDefine()
		}
		define(os.Args[2], l)
	case "api":
		api(flags.port, l)
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [api | define] ARGS\n", os.Args[0])
	os.Exit(1)
}

func usageDefine() {
	fmt.Fprintf(os.Stderr, "usage: %s define WORD\n", os.Args[0])
	os.Exit(1)
}
