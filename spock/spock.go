package spock

import (
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

// LogLevel sets the loglevel
func LogLevel(level string) {
	switch level {
	case "high":
		log.SetLevel(log.DebugLevel)
	case "med":
		log.SetLevel(log.InfoLevel)
	case "low":
		log.SetLevel(log.ErrorLevel)
	}
}

// New instance of spock
func New(dir string, lambda *genie.Genie, tracker *sherlock.Sherlock) *Spock {
	log.SetLevel(log.InfoLevel)
	spock := &Spock{
		Lambda: lambda,
		Dir:    dir,
		Cron:   cron.New(),
		Track:  tracker,
	}

	// Create our locks for our maps
	spock.Locks = make(map[string]*sync.Mutex)
	spock.Locks["Channels"] = &sync.Mutex{}
	spock.Locks["Checks"] = &sync.Mutex{}
	spock.Locks["Info"] = &sync.Mutex{}

	// load defaults out of the box
	spock.LoadDefaults()

	// continously watch channels and check
	spock.WatchChannelsAndChecks()

	// spock as the conn *nerd alert!*
	spock.Conn()

	log.Info("Starting Spock ...")

	// return the goods
	return spock
}

// Spock Singleton to handle all of the inner workings of our app
type Spock struct {
	Track    *sherlock.Sherlock     `yaml:"-"`
	Lambda   *genie.Genie           `yaml:"-"`
	Dir      string                 `yaml:"-"`
	Channels map[string]Channel     `yaml:",inline"`
	Checks   map[string]Check       `yaml:",inline"`
	Locks    map[string]*sync.Mutex `yaml:"-"`
	Cron     *cron.Cron             `yaml:"-"`
}

// convert a spock runner into a cron job
type job struct {
	s    *Spock
	name string
}

// run the cron job
func (j *job) Run() {
	check, exists := j.s.GetCheck(j.name)
	if exists {
		j.s.Runner(j.name, check)
	} else {
		log.WithFields(log.Fields{"name": j.name, "check": "failed"}).Error("Check doesn't exist")
	}
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

	for name, channel := range s.Channels {
		for lambdaType, lambdaValue := range channel.Lambdas {
			switch lambdaType {
			case "custom":
				s.Lambda.AddLambda(genie.NewCustomLambda(name, lambdaValue))
			case "github":
				fallthrough
			case "github.com":
				gh := strings.SplitAfterN(lambdaValue, "/", 3)
				if len(gh) == 3 {
					// perfect
					s.Lambda.GithubLambda(name, gh[0], gh[1], gh[2])
				} else {
					log.WithFields(log.Fields{"channel": name}).Error("Unable to find username/project/file for github")
				}
			case "code":
				if channel.Command != "" {
					l, err := genie.NewLocalLambda(name, s.Dir+"lambdas/", channel.Command, []byte(lambdaValue))
					if err == nil {
						s.Lambda.AddLambda(l)
					} else {
						log.WithFields(log.Fields{"channel": name}).Error(err.Error())
					}
				} else {
					log.WithFields(log.Fields{"channel": name}).Error("Missing 'command'")
				}
			default:
				if lambdaValue == "slack" {
					s.Lambda.AddLambda(genie.NewCodeLambda(name, channels.Slack))
				}
			}
			// only use the first one
			break
		}
	}

	s.Locks["Channels"].Unlock()
	return contents
}

// LoadDefaults will load all of our defaults to get spock running out of the box
func (s *Spock) LoadDefaults() {
	// add our url check
	s.Lambda.AddLambda(genie.NewCodeLambda("url", channels.URL))

	// add slack integration
	s.Lambda.AddLambda(genie.NewCodeLambda("slack", channels.Slack))

	// add shell lambda
	s.Lambda.AddLambda(genie.NewCustomLambda("shell", ""))
	s.Lambda.AddLambda(genie.NewCustomLambda("cmd", ""))
	s.Lambda.AddLambda(genie.NewCustomLambda("cli", ""))
	// ^^^ ha! genius if I do say so myself ... allows for custom cli/shell commands!

	// add heartbeats
	s.Lambda.AddLambda(genie.NewCodeLambda("heartbeat", channels.Heartbeat))
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
	} else {
		// boo!
		log.WithFields(log.Fields{"location": file}).Error(err.Error())
	}
	return ""
}

// WatchChannelsAndChecks will continously watch the channels and checks and reload them
func (s *Spock) WatchChannelsAndChecks() {
	// load for the very first time
	s.LoadChecks()
	s.LoadChannels()
	s.Conn()

	go func() {
		for {
			// do this action every 10 seconds
			<-time.After(10 * time.Second)

			reloadCron := false
			cha, chaErr := os.Stat(s.Dir + "channels/")
			if chaErr == nil && cha.ModTime().Add(10*time.Second).After(time.Now()) {
				log.WithFields(log.Fields{"location": s.Dir + "channels/"}).Info("Reloading configuration")
				s.LoadChannels()
				reloadCron = true
			}

			che, cheErr := os.Stat(s.Dir + "checks/")
			if cheErr == nil && che.ModTime().Add(10*time.Second).After(time.Now()) {
				log.WithFields(log.Fields{"location": s.Dir + "checks/"}).Info("Reloading configuration")
				s.LoadChecks()
				reloadCron = true
			}

			// spock has the conn ...
			if reloadCron {
				s.Conn()
			}
		}
	}()
}

// Conn controls the entire package
func (s *Spock) Conn() {
	s.Cron.Stop()
	s.Cron = cron.New()
	s.Locks["Checks"].Lock()
	defer s.Locks["Checks"].Unlock()
	for name, check := range s.Checks {
		// collect some stats
		e := s.Track.E(name)

		// if every is set
		if check.Every != "" {
			if strings.HasPrefix(check.Every, "@") {
				s.Cron.AddJob(check.Every, &job{s: s, name: name})
				e.S("interval").Set(check.Every)
			} else {
				s.Cron.AddJob("@every "+check.Every, &job{s: s, name: name})
				e.S("interval").Set("every " + check.Every)
			}
		}

		// cron is set ...
		if check.Cron != "" {
			s.Cron.AddJob(check.Cron, &job{s: s, name: name})
			e.S("interval").Set(check.Cron)
		}

		if check.Cron == "" && check.Every == "" {
			// you leave me no choice. Every thirty seconds it is! (todo: we _should_ pass this in as an arg)
			s.Cron.AddJob("*/30 * * * * *", &job{s: s, name: name})
			e.S("interval").Set("*/30 * * * * *")
		}
	}
	s.Cron.Start()
}

// GetChannel returns the channel and if it doens't, an error.
func (s *Spock) GetChannel(name string) (Channel, bool) {
	s.Locks["Channels"].Lock()
	channel, exists := s.Channels[name]
	s.Locks["Channels"].Unlock()
	return channel, exists
}

// GetCheck returns the given check
func (s *Spock) GetCheck(name string) (Check, bool) {
	s.Locks["Checks"].Lock()
	check, exists := s.Checks[name]
	s.Locks["Checks"].Unlock()
	return check, exists
}
