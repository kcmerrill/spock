package spock

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"time"

	"github.com/kcmerrill/crush/crush"
	"github.com/kcmerrill/genie/genie"
	"github.com/kcmerrill/spock/channels"
	"gopkg.in/yaml.v2"
)

// New instance of spock
func New(dir string, queue *crush.Q, lambda *genie.Genie) *Spock {
	spock := &Spock{
		Q:      queue,
		Lambda: lambda,
		Dir:    dir,
	}

	// load defaults out of the box
	spock.LoadDefaults()

	// load channels and checks
	spock.LoadChannels()
	spock.LoadChecks()

	// continously watch channels and check
	spock.WatchChannelsAndChecks()

	// return the goods
	return spock
}

// Spock Singleton to handle all of the inner workings of our app
type Spock struct {
	Q        *crush.Q            `yaml:"-"`
	Lambda   *genie.Genie        `yaml:"-"`
	Dir      string              `yaml:"-"`
	Channels map[string]*Channel `yaml:",inline"`
	Checks   map[string]*Check   `yaml:",inline"`
}

// LoadChecks loads all of the checks we will need to create.
func (s *Spock) LoadChecks() string {
	contents := s.loader(s.Dir + "checks/")
	yaml.Unmarshal([]byte(contents), &s.Checks)
	return contents
}

// LoadChannels will grab all of the channels spock needs to create.
func (s *Spock) LoadChannels() string {
	contents := s.loader(s.Dir + "channels/")
	yaml.Unmarshal([]byte(contents), &s.Channels)
	return contents
}

// LoadDefaults will load all of our defaults to get spock running out of the box
func (s *Spock) LoadDefaults() {
	// add our url check
	s.Lambda.AddLambda(genie.NewCodeLambda("url", channels.URL))
}

func (s *Spock) loader(file string) string {
	if _, err := os.Stat(file); err == nil {
		if files, dirErr := ioutil.ReadDir(file); dirErr == nil {
			contents := ""
			// directory
			for _, f := range files {
				if content, err := ioutil.ReadFile(file + f.Name()); err == nil {
					contents += strings.TrimSpace(string(content)) + "\n\n"
				}
			}
			return contents
		}

		// file
		if contents, err := ioutil.ReadFile(file); err == nil {
			log.Println(contents)
			return string(contents)
		}
	}
	// boo!
	return ""
}

// WatchChannelsAndChecks will continously watch the channels and checks and reload them
func (s *Spock) WatchChannelsAndChecks() {
	go func() {
		for {
			// Load every thirty seconds. We can make this configurable later
			<-time.After(30 * time.Second)

			// load channels and checks
			s.LoadChannels()
			s.LoadChecks()
		}
	}()
}
