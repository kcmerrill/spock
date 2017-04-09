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
func URL(stdin io.Reader, args []string) (string, error) {
	// lets get started ...
	var url, contains string
	var status int

	// our flags
	f := flag.NewFlagSet("url", flag.ContinueOnError)
	f.StringVar(&contains, "contains", "", "Simple string to look for")
	f.IntVar(&status, "status", 200, "Status code site must respond to")

	url = strings.TrimSpace(args[0])

	// set flags
	f.Parse(args[1:])

	if url != "" {
		response, respErr := http.Get(url)
		if respErr == nil {
			// no errors thus far ...
			body, _ := ioutil.ReadAll(response.Body)
			defer response.Body.Close()

			// check status code
			if status != 0 && response.StatusCode != status {
				return string(body), fmt.Errorf("Expected %d status code, Actual %d", response.StatusCode, status)
			}

			// check contains(search string)
			if contains != "" && !strings.Contains(string(body), contains) {
				return string(body), fmt.Errorf("Expected %s", contains)
			}

			// good to go!
			return string(body), nil
		}
		// something happened :(
		return respErr.Error(), respErr
	}
	return "", fmt.Errorf("Invalid URL '%s'", url)
}
