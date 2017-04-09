package spock

// Channel contains information about a channel
type Channel struct {
	Params   string            `yaml:"params"`
	Command  string            `yaml:"command"`
	Lambdas  map[string]string `yaml:",inline"`
	Template string            `yaml:"template"`
}
