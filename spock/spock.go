package spock

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kcmerrill/crush/crush"
	"github.com/kcmerrill/genie/genie"
	"github.com/kcmerrill/spock/channels"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// New instance of spock
func New(dir string, Q *crush.Q, lambda *genie.Genie) *Spock {
	spock := &Spock{
		Q:      Q,
		Lambda: lambda,
		Dir:    dir,
		Cron:   cron.New(),
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
	Q        *crush.Q               `yaml:"-"`
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
	s.Worker("url", 10)
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
			s.Cron.AddFunc(check.Cron, func() { s.Producer(name, check) })
			log.WithFields(log.Fields{"name": name, "cron": check.Cron}).Debug("Loaded")
		} else {
			// come on man! give me something!
			s.Cron.AddFunc("0 * * * * *", func() { s.Producer(name, check) })
			log.WithFields(log.Fields{"name": name, "every": "default(1m)"}).Debug("Loaded")
		}
	}
	s.Cron.Start()
}

// Producer takes a check and runs with it
func (s *Spock) Producer(name string, check *Check) {
	for topic := range check.GetMessages() {
		// create a new message
		msg := crush.NewMessage(topic, name, "")

		// we will take care of attempts, notifications, etc later ...
		msg.Attempts = 1

		// How long does it take?
		if check.Takes == "" {
			msg.Flight = "1m"
		} else {
			msg.Flight = check.Takes
		}

		// setup our dead letter(aka our notification channels)
		msg.DeadLetter = check.Notify

		// send our message to the queue
		s.Q.NewRawMessage(msg)
	}
}

// Worker starts a lambda worker based on a channel
func (s *Spock) Worker(channel string, workers int) {
	for worker := 0; worker < workers; worker++ {
		go func() {
			for {
				msg := s.Q.Message(channel)
				if msg == nil {
					// no need to hammer the system
					<-time.After(1 * time.Second)
					continue
				}

				// we have a message, lets process it
				s.Locks["Checks"].Lock()
				check, exists := s.Checks[msg.ID]
				s.Locks["Checks"].Unlock()
				if exists {
					for lambda, args := range check.GetMessages() {
						_, err := s.Lambda.Execute(lambda, strings.NewReader(""), args)
						if err == nil {
							// everything worked :thumbs_up:
						}
					}
					// remove the message regardless
					s.Q.Delete(channel, msg.ID)
				}
			}
		}()
	}
}
