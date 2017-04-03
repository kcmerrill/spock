package spock

import "strings"

// Check contains information needed for our checks
type Check struct {
	Params string            `yaml:"params"`
	Notify string            `yaml:"notify"`
	Module map[string]string `yaml:",inline"`
	Cron   string            `yaml:"cron"`
	Try    int               `yaml:"try"`
	Takes  string            `yaml:"takes"`
}

// GetMessages will grab all the modules + arguments
func (c *Check) GetMessages() map[string]string {
	msgs := make(map[string]string)

	for key, value := range c.Module {
		msgs[key] = "--" + key + "=" + value + " " + strings.Join(params(c.Params), " ")
	}

	return msgs
}
