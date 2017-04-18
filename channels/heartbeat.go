package channels

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/kcmerrill/spock/info"
)

// Heartbeat will see if a ping has happened
func Heartbeat(stdin io.Reader, args []string) (string, error) {
	// lets get started ...
	duration := strings.TrimSpace(args[0])
	in, _ := ioutil.ReadAll(stdin)
	cInfo := info.New(in)

	dur, err := time.ParseDuration(duration)
	if err != nil {
		return "Unable to parse heartbeat duration", fmt.Errorf("Unable to parse heartbeat duration")
	}

	if cInfo.Heartbeat.IsZero() {
		// nothing eh?
		return "Have not seen a heartbeat", fmt.Errorf("Have not seen a heartbeat")
	}

	if cInfo.Heartbeat.Add(dur).Before(time.Now()) {
		return "Over " + duration + " since a heartbeat has last been seen", fmt.Errorf("Over " + duration + " since a heartbeat has last been seen")
	}

	return "", nil
}
