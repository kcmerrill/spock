package main

import (
	"flag"
	"time"

	"github.com/kcmerrill/spock/spock"
)

func main() {
	checks := flag.String("checks", "./t/checks", "Checks directory")
	channels := flag.String("channels", "./t/channels", "Channels directory")
	reload := flag.Duration("reload", 10*time.Second, "How often the checks should be reloaded")
	flag.Parse()

	spock.New(&spock.Spock{
		ChecksDir:            *checks,
		ChannelsDir:          *channels,
		reloadConfigInterval: *reload,
	})
}
