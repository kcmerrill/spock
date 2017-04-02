package spock

// Check contains information needed for our checks
type Check struct {
	Params string            `yaml:"params"`
	Notify string            `yaml:"notify"`
	Module map[string]string `yaml:",inline"`
}
