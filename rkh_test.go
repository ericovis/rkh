package main

import "testing"

func TestRkh(t *testing.T) {
	FilePath = "./known_hosts"

	knownHosts := getKnownHosts()

	if len(*knownHosts) != 27 {
		t.Error("knownHosts := getKnownHosts() failed, length should be", 27, "but", len(*knownHosts), "was retrieved.")
	} else {
		t.Log("knownHosts := getKnownHosts() pass")
	}
}
