package channels

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// URL is a simple http check
func URL(stdin io.Reader, args string) (string, error) {
	// lets get started ...
	var url, contains string
	var status int

	// our flags
	f := flag.NewFlagSet("url", flag.ContinueOnError)
	f.StringVar(&url, "url", "", "URL to be tested")
	f.StringVar(&contains, "contains", "", "Simple string to look for")
	f.IntVar(&status, "status", 0, "Status code site must respond to")

	// set flags
	f.Parse(strings.Split(args, " "))

	if url != "" {
		response, respErr := http.Get(url)
		if respErr == nil {
			// no errors thus far ...
			body, _ := ioutil.ReadAll(response.Body)
			defer response.Body.Close()

			fmt.Println("BLEH", response.StatusCode, status)
			// check status code
			if status != 0 && response.StatusCode != status {
				return string(body), fmt.Errorf("Status code incorrect. Expected %d Actual %d", response.StatusCode, status)
			}

			// check contains(search string)
			if contains != "" && !strings.Contains(string(body), contains) {
				return string(body), fmt.Errorf("Search string not found! Expected to find %s", contains)
			}

			// good to go!
			return string(body), nil
		}
		// something happened :(
		return respErr.Error(), respErr
	}
	return "", fmt.Errorf("Invalid URL '%s'", url)
}
