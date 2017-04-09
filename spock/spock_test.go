package spock

import (
	"testing"

	"github.com/kcmerrill/genie/genie"
)

func TestChannels(t *testing.T) {
	s := New("../test/sample_b/", genie.New("test/lambdas", "", ""))
	if c, exists := s.Channels["slack"]; !exists {
		t.Fatalf("Unable to parse channels 'slack'")
	} else {
		// the channel params should be set
		if c.Params != "webhook=1234 channel=#general" {
			t.Fatalf("Unable to parse params from the channels yaml file")
		}
	}
}

func TestChecks(t *testing.T) {
	s := New("../test/sample_b/", genie.New("test/lambdas", "", ""))

	if _, exists := s.Checks["bad-kcmerrill.com"]; !exists {
		t.Fatalf("Expecting the check 'bad-kcmerrill.com'")
	}

	if c, exists := s.Checks["kcmerrill.com"]; !exists {
		t.Fatalf("Expecting the check 'kcmerrill.com'")
	} else {
		// just verify we got the goods back :)
		if c.Params != "status=200 contains=digital" {
			t.Fatalf("Expeced kcmerrill.com params to be 'status=200 contains=digital")
		}

		// verify notifications
		if c.Notify != "email slack" {
			t.Fatalf("Notification for kcmerrill.com should be email and slack")
		}

		// Verify the modules are being parsed properly
		if url, modExists := c.Module["url"]; modExists {
			if url != "http://kcmerrill.com" {
				t.Fatalf("Unable to parse modules. In this case kcmerrill.com -> url: http://kcmerrill.com")
			}
		} else {
			t.Fatalf("kcmerrill.com -> url -> kcmerrill.com should exist!")
		}
	}
}

func TestLoader(t *testing.T) {
	// we are in the spock dir ... go up one
	s := New("../test/sample_a/", genie.New("test/lambdas", "", ""))
	channels := s.LoadChannels()
	if channels != "a\n\nb\n\n" {
		t.Fatalf(channels)
	}

	// test our checks
	checks := s.LoadChecks()
	if checks != "checks\n\n" {
		t.Fatalf("Faied to load sample_a/checks/checks.yml -> contains 'a'")
	}
}

func TestGetChannel(t *testing.T) {
	// we are in the spock dir ... go up one
	s := New("../test/sample_a/", genie.New("test/lambdas", "", ""))
	channels := s.LoadChannels()
	if channels != "a\n\nb\n\n" {
		t.Fatalf(channels)
	}

	// test our checks
	checks := s.LoadChecks()
	if checks != "checks\n\n" {
		t.Fatalf("Faied to load sample_a/checks/checks.yml -> contains 'a'")
	}
}

// verify our defaults are loaded
func TestLoadDefaults(t *testing.T) {
	s := New("../test/sample_a/", genie.New("test/lambdas", "", ""))
	if _, exists := s.Lambda.Lambdas["slack"]; !exists {
		t.Fatalf("Expecting the 'slack' lambda to be enabled by default")
	}

	if _, exists := s.Lambda.Lambdas["url"]; !exists {
		t.Fatalf("Expecting the 'url' lambda to be enabled by default")
	}
}
