package main

import (
	"flag"

	"github.com/kcmerrill/spock/spock"
)

func main() {
	checks := flag.String("checks", "./t/checks", "Checks directory")
	channels := flag.String("channels", "./t/channels", "Channels directory")
	flag.Parse()

	spock.New(&spock.Spock{
		ChecksDir:   *checks,
		ChannelsDir: *channels,
	})
}
