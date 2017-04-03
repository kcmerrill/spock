package spock

import "testing"

func TestBasicCheck(t *testing.T) {
	c := &Check{
		Params: "arg1=foo arg2=bar arg3=true",
	}
	c.Module = make(map[string]string)
	c.Module["url"] = "http://kcmerrill.com"
	// while we are here, lets check the params to the CLI as well
}
