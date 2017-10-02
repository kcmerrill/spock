package spock

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v2"

	sherlock "github.com/kcmerrill/sherlock/core"
)

// New inits an instance of spock
func New(s *Spock) {
	// lets init some basics here
	s.stats = sherlock.New(100)
	s.checks = make(map[string]*check)
	s.channels = make(map[string]*channel)
	s.channelsLock = &sync.Mutex{}
	s.checksLock = &sync.Mutex{}
	s.ChannelsDir = strings.TrimRight(s.ChannelsDir, "/") + "/"
	s.ChecksDir = strings.TrimRight(s.ChecksDir, "/") + "/"

	// Start our checks loader
	go func() {
		for {
			s.loadChecks()
			<-time.After(time.Second)
		}
	}()

	// Start our channels loader
	go func() {
		for {
			s.loadChannels()
			<-time.After(time.Second)
		}
	}()

	// spock has the conn
	s.conn()
}

// Spock holds all of the configuration
type Spock struct {
	RootDir      string
	ChecksDir    string
	ChannelsDir  string
	checks       map[string]*check
	channels     map[string]*channel
	channelsLock *sync.Mutex
	checksLock   *sync.Mutex
	stats        *sherlock.Sherlock
}

func (s *Spock) concatFiles(dir string) []byte {
	files, filesError := filepath.Glob(dir + "*.yml")
	if filesError != nil {
		return []byte{}
	}

	config := []byte{}
	for _, file := range files {
		contents, _ := ioutil.ReadFile(file)
		config = append(config, []byte("\n")...)
		config = append(config, contents...)
	}

	return config
}

func (s *Spock) loadChecks() {
	checks := make(map[string]*check)
	loadedChecks := s.concatFiles(s.ChecksDir)
	s.checksLock.Lock()
	defer s.checksLock.Unlock()
	yamlError := yaml.Unmarshal(loadedChecks, &checks)
	if yamlError != nil {
		// TODO: output error
		return
	}

	for id, check := range checks {
		// does this check exist?
		if _, exists := s.checks[id]; !exists {
			// lets add it
			s.checks[id] = check
			fmt.Println("new")
		} else {
			update := false
			// no? already exists? Ok ... lets see what's different
			if s.checks[id].Cmd != check.Cmd {
				update = true
			}

			if update {
				s.checks[id] = check
				fmt.Println("update")
			}
		}
	}
}

func (s *Spock) loadChannels() {
	channels := s.concatFiles(s.ChannelsDir)
	s.channelsLock.Lock()
	defer s.channelsLock.Unlock()
	yamlError := yaml.Unmarshal(channels, &s.channels)
	if yamlError != nil {
		// TODO: output error
		return
	}
	// TODO: Show updated/deleted channels
}

func (s *Spock) conn() {
	for {
		s.checksLock.Lock()
		//for id, check := range s.checks {
		//fmt.Println(id, ">", check.checked, time.Now())
		//s.checks[id].checked = time.Now()
		//}
		s.checksLock.Unlock()
		//now := time.Now().Round(time.Second)
		<-time.After(time.Second)
	}
}
