package main

import (
	"flag"
	"syscall"

	"github.com/kcmerrill/crush/crush"
	"github.com/kcmerrill/genie/genie"
	shutdown "github.com/kcmerrill/shutdown.go"
	"github.com/kcmerrill/spock/spock"
	log "github.com/sirupsen/logrus"
)

func main() {
	var level string
	flag.StringVar(&level, "log-level", "low", "The log level output(low|med|high)")
	flag.Parse()

	switch level {
	case "low":
		log.SetLevel(log.ErrorLevel)
		break
	case "med":
		log.SetLevel(log.InfoLevel)
		break
	default:
		log.SetLevel(log.DebugLevel)
	}

	dir := "./test/sample_c/"
	spock.New(
		dir,
		crush.CreateQ(),
		genie.New(dir, "", ""),
	)

	shutdown.WaitFor(syscall.SIGINT, syscall.SIGTERM)
}
