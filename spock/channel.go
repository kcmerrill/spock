package spock

import (
	b64 "encoding/base64"
)

type channel struct {
	Name        string
	Description string `yaml:"description"`
}

func (c *channel) id() string {
	return b64.StdEncoding.EncodeToString([]byte(c.Name + c.Description))
}
