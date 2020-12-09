package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func packages(vendor string) map[string]string {
	list := map[string]string{}

	if vendor == "ubuntu" {
		out, err := exec.Command("dpkg", "-l").Output()
		if err != nil {
			log.Fatal(err)
		}

		cmdString := fmt.Sprintf("%s", out)
		scanner := bufio.NewScanner(strings.NewReader(cmdString))

		var count int
		for scanner.Scan() {
			line := scanner.Text()
			entries := strings.Fields(line)

			if entries[0] == "ii" || entries[0] == "rc" {
				list[entries[1]] = entries[2]
			}

			count++
		}
	}

	if vendor == "centos" {
		out, err := exec.Command("yum", "list", "installed").Output()
		if err != nil {
			log.Fatal(err)
		}

		cmdString := fmt.Sprintf("%s", out)
		scanner := bufio.NewScanner(strings.NewReader(cmdString))

		var count int
		for scanner.Scan() {
			line := scanner.Text()

			entries := strings.Fields(line)

			var fix int
			if line == "Installed Packages" {
				fix = count
			}

			if count > fix {
				if len(entries) == 3 {
					list[entries[0]] = entries[1]
				} else if len(entries) < 2 {
					list[entries[0]] = "not defined"
				}
			}
			count++
		}
	}

	return list
}
