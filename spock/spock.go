package spock

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kcmerrill/genie/genie"
	"github.com/kcmerrill/sherlock/sherlock"
	"github.com/kcmerrill/spock/channels"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// New instance of spock
func New(dir string, lambda *genie.Genie) *Spock {
	log.SetLevel(log.InfoLevel)
	spock := &Spock{
		Lambda: lambda,
		Dir:    dir,
		Cron:   cron.New(),
		Track:  sherlock.New(),
	}

	log.SetLevel(log.DebugLevel)

	// Create our locks for our maps
	spock.Locks = make(map[string]*sync.Mutex)
	spock.Locks["Channels"] = &sync.Mutex{}
	spock.Locks["Checks"] = &sync.Mutex{}
	spock.Locks["Info"] = &sync.Mutex{}

	// load defaults out of the box
	spock.LoadDefaults()

	// load channels and checks(at least initially()
	spock.LoadChannels()
	spock.LoadChecks()

	// continously watch channels and check
	spock.WatchChannelsAndChecks()

	// spock as the conn *nerd alert!*
	spock.Conn()

	// return the goods
	return spock
}

// Spock Singleton to handle all of the inner workings of our app
type Spock struct {
	Track    *sherlock.Sherlock     `yaml:"-"`
	Lambda   *genie.Genie           `yaml:"-"`
	Dir      string                 `yaml:"-"`
	Info     map[string]*Info       `yaml:"-"`
	Channels map[string]*Channel    `yaml:",inline"`
	Checks   map[string]*Check      `yaml:",inline"`
	Locks    map[string]*sync.Mutex `yaml:"-"`
	Cron     *cron.Cron             `yaml:"-"`
}

// LoadChecks loads all of the checks we will need to create.
func (s *Spock) LoadChecks() string {
	contents := s.loader(s.Dir + "checks/")
	s.Locks["Checks"].Lock()
	yaml.Unmarshal([]byte(contents), &s.Checks)
	s.Locks["Checks"].Unlock()
	return contents
}

// LoadChannels will grab all of the channels spock needs to create.
func (s *Spock) LoadChannels() string {
	contents := s.loader(s.Dir + "channels/")
	s.Locks["Channels"].Lock()
	yaml.Unmarshal([]byte(contents), &s.Channels)
	s.Locks["Channels"].Unlock()
	return contents
}

// LoadDefaults will load all of our defaults to get spock running out of the box
func (s *Spock) LoadDefaults() {
	// add our url check
	s.Lambda.AddLambda(genie.NewCodeLambda("url", channels.URL))

	// add slack integration
	s.Lambda.AddLambda(genie.NewCodeLambda("slack", channels.Slack))
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
			return string(contents)
		}
	}
	// boo!
	log.WithFields(log.Fields{"location": file}).Error("Error reading contents")
	return ""
}

// WatchChannelsAndChecks will continously watch the channels and checks and reload them
func (s *Spock) WatchChannelsAndChecks() {
	go func() {
		for {
			// do this action every 1minute, at least for now
			<-time.After(1 * time.Second)

			// let everybody know
			log.WithFields(log.Fields{"location": s.Dir}).Debug("Reloading configuration")

			// load channels and checks
			s.LoadChecks()
			s.LoadChannels()

			// spock has the conn ...
			s.Conn()
		}
	}()
}

// Conn controls the entire package
func (s *Spock) Conn() {
	s.Cron.Stop()
	s.Cron = cron.New()
	for name, check := range s.Checks {
		// cron is set ...
		if check.Cron != "" {
			s.Cron.AddFunc(check.Cron, func() { s.Runner(name, check) })
		} else {
			// come on man! give me something! ok ok, every minute it is.
			s.Cron.AddFunc("*/10 * * * * *", func() { s.Runner(name, check) })
		}
	}
	s.Cron.Start()
}

// Runner takes a check and runs with it
func (s *Spock) Runner(name string, check *Check) {
	for module, args := range check.GetMessages() {
		go func() {
			e := s.Track.E(name + "." + module)
			// might be multiple checks within a check. :shrug:
			e.I("checked").Add(1)
			// set the type
			e.S("module").Set(module)
			// set the name
			e.S("name").Set(name)
			// Reset the time it was last_checked to now()
			e.D("last_checked").Reset()
			// Update attempted
			e.I("attempts").Add(1)
			// take whatever we have currently in our entity and json encode it
			j, _ := json.Marshal(e)
			// Let execute it!
			results, ok := s.Lambda.Execute(module, strings.NewReader(string(j)), args)
			// noooice!
			if ok == nil {
				log.WithFields(log.Fields{"name": name, "module": module}).Info(string(j))
				// lets track last success
				e.D("last_success").Reset()
				// reset the attempts
				e.I("attempts").Reset()
			} else {
				log.WithFields(log.Fields{"name": name, "module": module, "check": "failed"}).Error(ok.Error())
				// boo! something broke ....
				e.D("last_failed").Reset()
				e.S("output").Set(results)
				e.S("error").Set(ok.Error())
				attempts := e.I("attempts").Int()
				if (check.Try == 0 && attempts == 1) || (check.Try == attempts) {
					// we should do something, like notify!
					e.D("last_notified").Reset()
					for _, notify := range strings.Split(check.Notify, " ") {
						j, _ := json.Marshal(e)
						go func(notify, module string) {
							// verify we have the appropriate channels
							if channel, exists := s.Channels[notify]; exists {
								// take whatever we have currently in our entity and json encode it
								_, notifyE := s.Lambda.Execute(notify, strings.NewReader(string(j)), strings.Join(params(channel.Params), " "))
								if notifyE != nil {
									log.WithFields(log.Fields{"name": name, "module": notify}).Error(notifyE.Error())
								}
							}
						}(notify, module)
					}
				}
			}
		}()
	}
}
