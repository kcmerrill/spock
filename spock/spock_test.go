package spock

import (
	"testing"

	"github.com/kcmerrill/crush/crush"
	"github.com/kcmerrill/genie/genie"
)

func TestChecks(t *testing.T) {
	s := New("../test/sample_b/", crush.CreateQ(), genie.New("test/lambdas", "", ""))
	/*
			kcmerrill.com:
		    url: kcmerrill.com
		    params: status=200 contains=digital https=true
		    notify: email slack
	*/
	if c, exists := s.Checks["kcmerrill.com"]; !exists {
		t.Fatalf("Expecting the check 'kcmerrill.com'")
	} else {
		if url, urlExists := c.Module["url"]; !urlExists || url != "kcmerrill.com" {
			t.Fatalf("Expecting the url to be 'kcmerrill.com")
		}

		if c.Params != "status=200 contains=digital" {
			t.Fatalf("Expecting parmas to be 'status=200 contains=digital'")
		}

		// while we are here, lets check the params to the CLI as well
		cli := c.ParamsCli(c.Params)

	}
}

func TestLoader(t *testing.T) {
	// we are in the spock dir ... go up one
	s := New("../test/sample_a/", crush.CreateQ(), genie.New("test/lambdas", "", ""))
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
