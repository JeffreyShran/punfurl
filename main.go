// USAGE:
//
// root@workstation:~/tools/punfurl# echo "https://support.google.com/google-ads/answer/2472708?hl=en-GB"  | go run ~/tools/punfurl/main.go
// https://support.google.com/google-ads
// https://support.google.com/answer
// https://support.google.com/google-ads/answer
// https://support.google.com/2472708
// https://support.google.com/google-ads/2472708
// https://support.google.com/answer/2472708
// https://support.google.com/google-ads/answer/2472708
//

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"net/url"
	"regexp"
)

func main() {

	var verbose bool
	flag.BoolVar(&verbose, "v", false, "")
	flag.BoolVar(&verbose, "verbose", false, "")

	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024) // Now we can accept 1mb lines

//	seen := make(map[string]bool)

	for scanner.Scan() {

		// Set up regex to find base64 urls
		found, err := regexp.MatchString(`data\:\w+\/\w+\;base64`, scanner.Text())
		if err != nil {
                	if verbose {
                        	fmt.Fprintf(os.Stderr, "parse failure: %s\n", err)
                        }
                }

		// base64 urls are usually images that we don't want.
		if found {
			continue
        	}

		// we have some confidence that we want to parse the url at this point so do it.
		u, err := parseURL(scanner.Text())
		if err != nil {
			if verbose {
				fmt.Fprintf(os.Stderr, "parse failure: %s\n", err)
			}
			continue
		}

		// Splits the path by the slashes
		f := func(c rune) bool {
    			return c == '/'
		}

		routes := []string{}

		// I had some garbage come through that I never want so i throw it out
		// and create the routes for the next stage
		for _, r := range strings.FieldsFunc(u.Path, f) {
			if ! strings.ContainsAny(r, ".") {
				routes = append(routes, r)
			}
		}

		// Make it safer and easier to check for dupes
		// split into 2 parts because golang hates me
		//dupeString := strings.Join(routes[:], "/")
		//dupeString = u.Scheme + string("://") + u.Host + string("/") + dupeString

		// Skip duplicates if we've seen them before
  //              if seen[dupeString] {
//					if verbose {
//						fmt.Println(string("skipping - ") + dupeString) // This should probably be a stderror
//					}
//					continue
  //              }

		// Take a note that we've seen this one before for checks on the next iteration
//		seen[dupeString] = true

		// The powerset allows us to keep the position of the paths but also
		// gives a nicely succinct variation on each one which is great for fuzzing.
		for _, r := range PowerSet(routes) {
			if len(r) > 0 {
				fmt.Println(u.Scheme + string("://") + u.Host + string("/") + strings.Join(r[:], "/"))
			}
		}
        // Print the host so we don't lose it
        fmt.Println(u.Scheme + string("://") + u.Host)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

// parseURL parses a string as a URL and returns a *url.URL
// or any error that occured. If the initially parsed URL
// has no scheme, http:// is prepended and the string is
// re-parsed
// https://github.com/tomnomnom/unfurl/blob/master/main.go
func parseURL(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	if ( u.Scheme == "" ) {
		return url.Parse("http://" + raw)
	}

	return u, nil
}

func copyAndAppendString(slice []string, elem string) []string {
	return append(append([]string(nil), slice...), elem)
}

// PowerSet creates unique combinations from a provided array
// not every combination and it keeps position which
// is what we want for api's for example.
func PowerSet(s []string) [][]string {
	if s == nil {
		return nil
	}
	r := [][]string{[]string{}}
	for _, es := range s {
		var u [][]string
		for _, er := range r {
			u = append(u, copyAndAppendString(er, es))
		}
		r = append(r, u...)
	}
	return r
}
