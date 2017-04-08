package channels

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// Slack sends a messages as an incoming webhook to the slack api
func Slack(stdin io.Reader, args string) (string, error) {
	// lets get started ...
	var webhook, channel string

	// our flags
	f := flag.NewFlagSet("slack", flag.ContinueOnError)
	f.StringVar(&webhook, "webhook", "", "The integration endpoint")
	f.StringVar(&channel, "channel", "", "The channel to be used")

	// set flags
	f.Parse(strings.Split(args, " "))

	if webhook != "" {
		// yay! we have a webhook!
		in, _ := ioutil.ReadAll(stdin)
		info := &template{}
		json.Unmarshal(in, info)
		fmt.Println("Sending to slack!")
		fmt.Println(webhook, channel, info.ID)
	}

	return "", fmt.Errorf("webook not set")
}
