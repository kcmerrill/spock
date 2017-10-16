package spock

import (
	"log"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

// New inits an instance of spock
func New(s *Spock) {
	// lets init some basics here
	s.checks = make(map[string]*check)
	s.channels = make(map[string]*channel)
	s.channelsLock = &sync.Mutex{}
	s.checksLock = &sync.Mutex{}
	s.ChannelsDir = strings.TrimRight(s.ChannelsDir, "/") + "/"
	s.ChecksDir = strings.TrimRight(s.ChecksDir, "/") + "/"

	// spock has the conn
	s.conn()
}

// Spock coordinates channels and checks
type Spock struct {
	CheckInterval time.Duration

	ChecksDir  string
	checks     map[string]*check
	checksLock *sync.Mutex

	ChannelsDir  string
	channels     map[string]*channel
	channelsLock *sync.Mutex
}

// loadChecks reads in all config files and reloads the configuration
func (s *Spock) loadChecks() {
	checks := make(map[string]*checkConfig)
	loadedChecks := combineConfigFiles(s.ChecksDir)
	s.checksLock.Lock()
	defer s.checksLock.Unlock()
	yamlError := yaml.Unmarshal(loadedChecks, &checks)
	if yamlError != nil {
		log.Println("ERROR", "YAML parsing")
		return
	}

	for name, lc := range checks {
		// does this check exist?
		if _, exists := s.checks[name]; !exists {
			// ok. lets add it
			s.checks[name] = &check{
				Name:     name,
				Config:   lc,
				ExecLock: &sync.Mutex{},
			}
			log.Println("NEW CHECK", name)
		} else {
			if s.checks[name].id() != lc.id() {
				s.checks[name].Config = lc
				log.Println("UPDATED CHECK", name)
			}
		}
	}
}

// loadChannels reads in all config files and reloads the configuration
func (s *Spock) loadChannels() {
	channels := make(map[string]*channel)
	loadedChannels := combineConfigFiles(s.ChannelsDir)
	s.channelsLock.Lock()
	defer s.channelsLock.Unlock()
	yamlError := yaml.Unmarshal(loadedChannels, &channels)
	if yamlError != nil {
		log.Println("CHANNEL ERROR", "YAML parsing")
		return
	}

	for name, c := range channels {
		// does this check exist?
		if _, exists := s.channels[name]; !exists {
			// ok. lets add it
			s.channels[name] = &channel{}
			log.Println("NEW CHANNEL", name)
		} else {
			if s.channels[name].id() != c.id() {
				s.channels[name] = c
				log.Println("UPDATED CHANNEL", name)
			}
		}
	}
}

// conn is the central hub for dispatching checks
func (s *Spock) conn() {
	for {
		s.loadChecks()
		s.loadChannels()
		for _, check := range s.checks {
			if check.LastChecked.Add(check.interval()).Before(time.Now()) {
				check.check(s.channels)
			}
		}
		<-time.After(s.CheckInterval)
	}
}
