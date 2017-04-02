package spock

import "strings"

// params returns an array of strings to be used as cli arguments
func params(params string) []string {
	cli := make([]string, 0)

	// lets keep this really simple(for now)
	// convert "arg1=foo arg2=bar arg3=true" -> "--arg1=foo --arg2=bar --arg3"
	for _, param := range strings.Split(params, " ") {
		args := strings.Split(param, "=")
		if len(args) == 2 {
			// lets test booleans right quick
			if strings.ToLower(args[1]) == "true" {
				cli = append(cli, "--"+args[0])
				continue
			}

			// false? continue ... we can come back to this later
			if strings.ToLower(args[0]) == "false" {
				continue
			}

			// good to go! just add some dashes
			cli = append(cli, "--"+param)
		} else {
			// not sure what to do here yet ...
			// user passed in something goofy ... stop being goofy user!
		}
	}

	// return the goods
	return cli
}
