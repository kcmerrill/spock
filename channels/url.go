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
				errMsg := fmt.Errorf("Expected %d status code, %d returned", status, response.StatusCode)
				return errMsg.Error(), errMsg
			}

			// check contains(search string)
			if contains != "" && !strings.Contains(string(body), contains) {
				errMsg := fmt.Errorf("Expected %s", contains)
				return errMsg.Error(), errMsg
			}

			// good to go
			return string(body), nil
		}
		// something happened :(
		return respErr.Error(), respErr
	}
	return "", fmt.Errorf("Invalid URL '%s'", url)
}
