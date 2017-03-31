package main

import (
	"github.com/kcmerrill/crush/crush"
	"github.com/kcmerrill/genie/genie"
	"github.com/kcmerrill/spock/spock"
)

func main() {
	dir := "."
	spock.New(
		dir,
		crush.CreateQ(),
		genie.New(dir, "", ""),
	)
}
