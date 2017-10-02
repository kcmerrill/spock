package spock

import "time"

type check struct {
	Description string    `yaml:"desc"`
	Cmd         string    `yaml:"cmd"`
	Try         int       `yaml:"try"`
	Every       string    `yaml:"every"`
	checked     time.Time `yaml:"-"`
}
