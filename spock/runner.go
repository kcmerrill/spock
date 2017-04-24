package spock

import (
	"encoding/json"
	"strings"

	"github.com/kcmerrill/sherlock/sherlock"
	log "github.com/sirupsen/logrus"
)

// Runner takes a check and runs with it
func (s *Spock) Runner(name string, check Check) {
	for module, args := range check.GetMessages() {
		go func() {
			e := s.Track.E(name)
			// might be multiple checks within a check. :shrug:
			e.I("checked").Add(1)
			// set the module
			e.S("module").Set(module)
			// set the name
			e.S("name").Set(name)
			// Reset the time it was last_checked to now()
			e.D("last_checked").Reset()
			// Update attempted
			e.I("attempts").Add(1)
			j, _ := json.Marshal(e)
			// Let execute it!
			results, ok := s.Lambda.Execute(module, strings.NewReader(string(j)), strings.Split(args, " "))
			// noooice!
			if ok == nil {
				log.WithFields(log.Fields{"name": name, "module": module, "#": e.I("checked").Int()}).Info("Check was successful!")
				// lets track last success
				e.D("last_success").Reset()
				// set the success/failure
				e.B("status").Set(true)
				// lets check for attempts right quick
				if e.B("notified").Bool() {
					// meaning, was bad, but is now good! Lets send it to a different channel
					s.Send("notify.recovery", check.Recovers, e)
				}
				// reset the attempts
				e.I("attempts").Reset()
				// create an event
				e.Event("checked:ok")
				// reset notified
				e.B("notified").Reset()
			} else {
				// set the success/failure
				e.B("status").Set(false)
				if strings.TrimSpace(results) != "" {
					// first try the lambda results ...
					log.WithFields(log.Fields{"name": name, "module": module, "attempts": e.I("attempts").Int()}).Error(strings.TrimSpace(results))
				} else {
					// perhaps the lambda doesn't exist? Or something broke with genie?
					log.WithFields(log.Fields{"name": name, "module": module, "attempts": e.I("attempts").Int()}).Error(ok.Error())
				}

				// boo! something broke ....
				e.D("last_failed").Reset()
				e.S("output").Set(results)
				e.S("error").Set(results)
				attempts := e.I("attempts").Int()
				// if try == 0 then we need to alert. If try == attempts we need to alert. If the last alert failed to alert, we need ot alert
				if (check.Try == 0 && attempts == 1) || (check.Try == attempts) || e.B("notification.error").Bool() {
					e.B("notified").Set(true)
					s.Send("notify.failure", check.Fails, e)
				}
				e.Event("checked:failed")
			}
		}()
	}
}

// Send takes in a space separated string of channels, and peforms the action based on the channels given
func (s *Spock) Send(rf, send string, e *sherlock.Entity) {
	// we should do something, like notify!
	e.D("last_notified").Reset()
	// reset all send errors
	e.B(rf + ".error").Set(false)
	for _, sendTo := range strings.Split(send, " ") {
		if sendTo == "" {
			// no need to jump through hoops if the dev didn't anything for notify
			continue
		}
		// verify we have the appropriate channels
		if channel, exists := s.GetChannel(sendTo); exists {
			// set our template to be passed along
			e.S("template").Set(channel.Template)
			// json encode our entity, so we can then use it later on for templating stuff
			j, _ := json.Marshal(e)
			// take whatever we have currently in our entity and json encode it
			notifyResults, notifyE := s.Lambda.Execute(sendTo, strings.NewReader(string(j)), strings.Split(channel.Params, " "))
			if notifyE != nil {
				log.WithFields(log.Fields{"name": e.S("name").String(), "channel": sendTo, "type": rf + ".error"}).Error(notifyE.Error())
				// try again next check(if it fails)
				e.B(rf + ".error").Set(true)
			} else {
				log.WithFields(log.Fields{"name": e.S("name").String(), "channel": sendTo, "type": rf + ".success"}).Info(notifyResults)
			}
		} else {
			log.WithFields(log.Fields{"name": e.S("name").String(), "channel": sendTo, "type": rf + ".error"}).Error("Channel does not exist")
			// try again next check(if it fails)
			e.B(rf + ".error").Set(true)
		}
	}
}
