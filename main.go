package main

import (
	"flag"
	"time"

	"github.com/kcmerrill/spock/spock"
)

func main() {
	checks := flag.String("checks", "./t/checks", "Checks directory")
	channels := flag.String("channels", "./t/channels", "Channels directory")
	checkInterval := flag.Duration("check-interval", time.Second, "How often the checks/channels should be reloaded")
	flag.Parse()

	spock.New(&spock.Spock{
		ChecksDir:     *checks,
		ChannelsDir:   *channels,
		CheckInterval: *checkInterval,
	})
}
