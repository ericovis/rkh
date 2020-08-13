package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var FilePath = os.Getenv("HOME") + "/.ssh/known_hosts"

func check(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

func getKnownHosts() *[]string {
	dat, err := ioutil.ReadFile(FilePath)
	check(err)

	lines := strings.Split(string(dat), "\n")

	return &lines
}


func removeHosts(matchedHosts *[]string, knownHosts *[]string) {
	var cleaned []string
	fmt.Println("Removing hosts...")
	for _, kHost := range *knownHosts {
		var match bool = false
		
		for _, mHost := range *matchedHosts {
			if strings.HasPrefix(kHost, mHost) {
				match = true
				break		
			}
		}

		if !match {
			cleaned = append(cleaned, kHost)
		}
	}
	contents := strings.Join(cleaned, "\n")
	err := ioutil.WriteFile(FilePath, []byte(contents), 0644)
	check(err)	
	fmt.Println("Done!")
}


func getMatchedHosts(knownHosts *[]string, args *[]string) []string {
	var matchedHosts []string

	re := regexp.MustCompile(`^([^\s]+)`)

	for _, host := range *knownHosts {
		for _, arg := range *args {
			match, _ := regexp.MatchString(arg, host)
			if match {
				matchedHosts = append(matchedHosts, re.FindString(host))
				break
			}
		}
	}

	return matchedHosts
}

func main() {
	flag.Usage = func() {
		usage := "Usage: rkh [parameters] [hostA hostB hostC ... hostZ]\n\nParameters:\n"
		fmt.Fprintf(os.Stderr, usage)
		flag.PrintDefaults()
	}

	forcePtr := flag.Bool("force", false, "Suppress confirmation dialog.")
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		os.Exit(128)
	}

	knownHosts := getKnownHosts()
	matchedHosts := getMatchedHosts(knownHosts, &args)

	fmt.Println(len(*knownHosts), "known host(s) were found.")

	if len(matchedHosts) == 0 {
		fmt.Println("No hosts were matched.")
		os.Exit(0)
	}

	canRemove := *forcePtr

	fmt.Println(len(matchedHosts), "host(s) were matched:")
	for _, host := range matchedHosts {
		fmt.Println(host)
	}

	if !canRemove {
		fmt.Println("\nConfirm the removal of these hosts? [Y/n]")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.Replace(answer, "\n", "", -1)

		if strings.ToLower(answer) == "yes" || strings.ToLower(answer) == "y" || strings.ToLower(answer) == "" {
			canRemove = true
		} else {
			fmt.Println("Exiting then...")
		}

	}

	if canRemove {
		removeHosts(&matchedHosts, knownHosts)
	}

	os.Exit(0)

}
