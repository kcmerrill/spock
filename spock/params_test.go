package spock

import "testing"

func TestParams(t *testing.T) {
	cli := params("arg1=foo arg2=bar arg3=true arg4=TRUE")

	if cli[0] != "--arg1=foo" {
		t.Fatalf("cli[0]='--arg1=foo'")
	}

	if cli[1] != "--arg2=bar" {
		t.Fatalf("cli[1]='--arg2=bar'")
	}

	// test bool
	if cli[2] != "--arg3" {
		t.Fatalf("cli[2]='--arg3'")
	}

	// test bool(case insensitive)
	if cli[3] != "--arg4" {
		t.Fatalf("cli[2]='--arg4'")
	}
}
