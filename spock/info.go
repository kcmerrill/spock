package spock

import "time"

// Info contains all the necessary informtion for our checks
type Info struct {
	Created  time.Time
	LastSeen time.Time
	Cadence  string
}
