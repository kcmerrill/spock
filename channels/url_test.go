package channels

import (
	"strings"
	"testing"
)

func TestURL(t *testing.T) {
	// check good functionality
	_, goodErr := URL(strings.NewReader(""), strings.Split("https://kcmerrill.com --status=200 --contains=digital", " "))
	if goodErr != nil {
		t.Fatalf(goodErr.Error())
	}

	// check bad search functionality
	_, errBadContains := URL(strings.NewReader(""), strings.Split("https://kcmerrill.com --status=200 --contains=bingowashisnameo", " "))
	if errBadContains == nil {
		t.Fatalf("kcmerrill.com does not contain the word 'bingowashisnameo'")
	}

	// check bad status code functionality
	_, errBadStatusCode := URL(strings.NewReader(""), strings.Split("http://kcmerrill.com --status=201 --contains=bingowashisnameo", " "))
	if errBadStatusCode == nil {
		t.Fatal("kcmerrill.com should return 200, not 201. So status codes are not working")
	}

	// verify redirects still work
	_, errRedirects := URL(strings.NewReader(""), strings.Split("http://kcmerrill.com --status=200 --contains=bingowashisnameo", " "))
	if errRedirects == nil {
		t.Fatal("Redirects should work too ... or kcmerrill.com is down :(")
	}
}
