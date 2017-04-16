package main

import (
	"flag"
	"fmt"
	"syscall"

	"github.com/kcmerrill/genie/genie"
	shutdown "github.com/kcmerrill/shutdown.go"
	"github.com/kcmerrill/spock/spock"
)

var (
	// Version will be set at runtime(the current version of spock)
	Version = "dev"
	// Commit will be set at runtime(the current commit id of spock)
	Commit = "n/a"
)

func main() {
	var logLevel, dir string
	var showVersion bool

	flag.StringVar(&logLevel, "v", "high", "Log level verbosity(low|med|high)")
	flag.StringVar(&dir, "dir", "./", "Root directory where your channels and checks are located")
	flag.BoolVar(&showVersion, "version", false, "Show Spock's version number")
	flag.Parse()

	// Show Version
	if showVersion {
		fmt.Println()
		fmt.Println("Spock - Making sure your applications live long and prosper.")
		fmt.Println("---")
		fmt.Println("Version: ", Version)
		fmt.Println("CommitId: ", Commit)
		fmt.Println("---")
		fmt.Println("Made with <3 by http://kcmerrill.com")
		fmt.Println()
		return
	}

	// set some log levels
	spock.LogLevel(logLevel)
	// disable genie level logging
	genie.LogLevel("panic")

	spock.New(
		dir,
		genie.New(dir+"lambdas/", "", ""),
	)

	// wait for shutdown signal
	shutdown.WaitFor(syscall.SIGINT, syscall.SIGTERM)
}
