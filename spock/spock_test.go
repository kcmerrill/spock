package spock

import (
	"fmt"
	"testing"

	"reflect"

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
		fmt.Println(c.Module)
		if url, urlExists := c.Module["url"]; !urlExists || url != "kcmerrill.com" {
			t.Fatalf("Expecting the url to be 'kcmerrill.com")
		}

		fmt.Println(reflect.TypeOf(c.Params))
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
