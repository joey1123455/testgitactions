package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type logAuthor struct {
	username string
	email    string
}

type gitLog struct {
	message string
	author  logAuthor
	id      string
	date    time.Time
}

type packageVersion struct {
	major int
	minor int
	patch int
}

func main() {
	// Set up the Git command
	// cmd := exec.Command("git", "log", "--format=%s")
	cmd := exec.Command("git", "log")

	// Set the working directory if needed
	// cmd.Dir = "/path/to/your/repo"

	// Create a pipe to capture the command output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Create an empty slice to store the git log messages
	var gitLogs []string

	// Read and store each line of the command output
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		gitLogs = append(gitLogs, line)
	}

	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error:", err)
	}

	logsList := parseLogs(gitLogs)
	// Print the git log messages
	for _, _ = range logsList {
		// println(log.id)
		// println(log.author.username)
		// println(log.author.email)
		// println(log.date.String())
		// println(log.message)
		// println()
		print(".")
	}
	println()
}

func parseLogs(logsList []string) []*gitLog {

	logs := make([]*gitLog, 0)
	//TODO: Retrieve from file
	packageVersion := packageVersion{}
	for _, logLine := range logsList {
		// println(logLine)
		var newLog gitLog
		if strings.HasPrefix(logLine, "commit") {
			// If it does, print the string excluding the "commit" prefix
			newLog = gitLog{
				id: strings.TrimSpace(logLine[len("commit "):]),
			}
			logs = append(logs, &newLog)

		} else {
			idx := len(logs) - 1
			// println(idx)
			switch {
			case strings.HasPrefix(logLine, "Author:"):
				pattern := regexp.MustCompile(`Author: (\w+) <([^>]+)>`)
				// Find matches in the line
				matches := pattern.FindStringSubmatch(logLine)
				// Check if there's a match and extract username and email
				if len(matches) == 3 {
					logs[idx].author.username = matches[1]
					logs[idx].author.email = matches[2]
				}

			case strings.HasPrefix(logLine, "Date:"):
				timeString := strings.TrimSpace(logLine[len("Date:"):])
				// fmt.Println(logLine)
				// println(timeString)
				layout := "Mon Jan 02 15:04:05 2006 -0700"

				// Parse the time string
				t, err := time.Parse(layout, timeString)
				if err != nil {
					fmt.Println("Error parsing time:", err)
				}

				// convert the time to utc
				logs[idx].date = t.UTC()

			default:
				regex := regexp.MustCompile(`^[a-zA-Z0-9@#$%^&*()-_+=! ]+$`)

				// Check if the string matches the regular expression
				if regex.MatchString(logLine) {
					commitMessage := strings.TrimSpace(logLine)

					// TODO: update package version
					updatePackageVersion(commitMessage, &packageVersion)
					logs[idx].message = commitMessage
					// println()
				}
				// logs[idx].message = strings.TrimSpace(logLine)
				// println(logs[idx].message)
			}
		}
	}
	print(packageVersion.major)
	print(".")
	print(packageVersion.minor)
	print(".")
	print(packageVersion.patch)
	println()
	return logs
}

func updatePackageVersion(commit string, version *packageVersion) {
	// Bug fixes patch
	// features minor version
	// breaking change major version
	// chore patch
	// default message patch
	const maxVersion = 99

	switch {
	case matchBugFixes(commit):
		println("bug fix")
		if version.patch < maxVersion {
			version.patch++
		} else {
			version.minor++
			version.patch = 0
		}
	case matchFeature(commit):
		println("feature")
		if version.minor < maxVersion {
			version.minor++
		} else {
			version.major++
			version.minor = 0
		}
	}
}

func matchBugFixes(comment string) bool {
	pattern := `^(bug)(\s)?(fix)(es\s*)?`

	// Compile the regular expression pattern
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(comment)
}

func matchFeature(comment string) bool {
	// pattern := `^(feat)(\s)?(new|feature|features|new feature|new-features|new feature(s))`
	pattern := `^(feat)`

	regex := regexp.MustCompile(pattern)

	return regex.MatchString(comment)
}
