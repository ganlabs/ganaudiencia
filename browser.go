package main

import (
	"math/rand/v2"
	"os/exec"
	"runtime"
)

func OpenBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
		// log.Println("windows")
	case "darwin":
		cmd = "open"
		// log.Println("macos")
	default:
		cmd = "xdg-open"
		// log.Println("linux")
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func GetPort() int {
	min := 7000
	max := 7900
	return (rand.IntN(max-min+1) + min)
}
