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

var filePath = os.Getenv("HOME") + "/.ssh/known_hosts"

func check(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

func getKnowHosts() *[]string {
	dat, err := ioutil.ReadFile(filePath)
	check(err)

	lines := strings.Split(string(dat), "\n")

	return &lines
}


func removeHosts(matchedHosts *[]string, knownHosts *[]string) {
	var cleaned []string
	fmt.Println("Removing hosts...")
	for _, mHost := range *matchedHosts {
		for _, kHost := range *knownHosts {
			if !strings.HasPrefix(kHost, mHost) {
				cleaned = append(cleaned, kHost)
			}
		}
	}
	contents := strings.Join(cleaned, "\n")
	err := ioutil.WriteFile(filePath, []byte(contents), 0644)
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

	knownHosts := getKnowHosts()
	matchedHosts := getMatchedHosts(knownHosts, &args)

	if len(matchedHosts) == 0 {
		fmt.Println("No hosts were matched.")
		os.Exit(0)
	}

	canRemove := *forcePtr

	if !canRemove {
		fmt.Println(len(matchedHosts), "host(s) were matched:")
		for _, host := range matchedHosts {
			fmt.Println(host)
		}

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
