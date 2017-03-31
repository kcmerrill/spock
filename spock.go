package main

import (
	"github.com/kcmerrill/crush/crush"
	"github.com/kcmerrill/genie/genie"
	"github.com/kcmerrill/spock/spock"
)

func main() {
	spock.New(".", crush.CreateQ(), genie.New("/tmp/", "", ""))
}
