package spock

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/kcmerrill/crush/crush"
	"github.com/kcmerrill/genie/genie"
	"gopkg.in/yaml.v2"
)

// New instance of spock
func New(dir string, queue *crush.Q, lambda *genie.Genie) *Spock {
	spock := &Spock{
		Q:      queue,
		Lambda: lambda,
		Dir:    dir,
	}

	// load channels and checks
	spock.LoadChannels()
	spock.LoadChecks()

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
