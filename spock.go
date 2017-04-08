package main

import (
	"syscall"

	"github.com/kcmerrill/genie/genie"
	shutdown "github.com/kcmerrill/shutdown.go"
	"github.com/kcmerrill/spock/spock"
)

func main() {
	dir := "./test/sample_b/"
	spock.New(
		dir,
		genie.New(dir, "", ""),
	)

	shutdown.WaitFor(syscall.SIGINT, syscall.SIGTERM)
}
