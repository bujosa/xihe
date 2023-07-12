package utils

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func CleanDns() {
	commands := []string{
		"sudo sysctl net.ipv6.conf.all.disable_ipv6=1",
		"sudo systemd-resolve --flush-caches",
	}

	cmd := exec.Command("/bin/sh", "-c", strings.Join(commands, " && "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Println("Error restarting the system", err)
	}
}
