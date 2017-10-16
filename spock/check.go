package spock

import (
	b64 "encoding/base64"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type check struct {
	Name        string
	ExecLock    *sync.Mutex
	Config      *checkConfig
	LastChecked time.Time
	Results     string
	Error       error
	Attempts    int
}

func (c *check) interval() time.Duration {
	return c.Config.interval()
}

func (c *check) id() string {
	return c.Config.id()
}

func (c *check) check(channels map[string]*channel) {
	go func() {
		c.ExecLock.Lock()
		defer c.ExecLock.Unlock()

		// make sure we update the last checked date
		c.LastChecked = time.Now()

		// alright, lets see what we get
		c.Results, c.Error = c.exec(c.Config.Check.Cmd)
		if c.Error == nil {
			c.Attempts = 0
			log.Println("CHECKED", c.Name)
			for okChannels := range c.channels(c.Config.Check.OK)
		} else {
		}

		// No good. Send it off to be processed
		c.Attempts++
		log.Println("FAILED", c.Name)
	}()
}

func (c *check) exec(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	output, error := cmd.CombinedOutput()
	return string(output), error
}

func (c *check) channels(channels string) []string {
	return strings.Split(channels, " ")
}

type checkConfig struct {
	Description string `yaml:"desc"`
	Check       struct {
		Cmd      string        `yaml:"cmd"`
		Tries    int           `yaml:"try"`
		Every    time.Duration `yaml:"every"`
		OK       string        `yaml:"ok"`
		Fails    string        `yaml:"fails"`
		Recovers string        `yaml:"recovers"`
	} `yaml:"check"`
}

func (cf *checkConfig) id() string {
	return b64.StdEncoding.EncodeToString([]byte(cf.Description + cf.Check.Cmd + cf.interval().String()))
}

func (cf *checkConfig) interval() time.Duration {
	if cf.Check.Every == (0 * time.Second) {
		return (30 * time.Second)
	}
	return cf.Check.Every
}
