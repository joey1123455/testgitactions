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
	for _, log := range logsList {
		println(log.id)
		println(log.author.username)
		println(log.author.email)
		println(log.date.String())
		println(log.message)
		println()
		// print(".")
	}
}

func parseLogs(logsList []string) []*gitLog {

	logs := make([]*gitLog, 0)
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

				// Print the parsed time
				logs[idx].date = t
				// println(t.String())
				// println(logs[idx].date.String())
				// println()

			default:
				regex := regexp.MustCompile(`^[a-zA-Z0-9@#$%^&*()-_+=! ]+$`)

				// Check if the string matches the regular expression
				if regex.MatchString(logLine) {
					fmt.Println("String contains spaces along with alphabets, numbers, or special symbols.")
					// println(logLine)
					logs[idx].message = strings.TrimSpace(logLine)
					println()
				}
				// logs[idx].message = strings.TrimSpace(logLine)
				// println(logs[idx].message)
			}
		}
	}
	return logs
}

func startsWithCommit(s string) bool {
	// Define the regular expression pattern
	pattern := `^commit\s`

	// Compile the regular expression pattern
	regex := regexp.MustCompile(pattern)

	// Use FindStringIndex to check if the string matches the pattern at the beginning
	if idx := regex.FindStringIndex(s); idx != nil && idx[0] == 0 {
		return true
	}
	return false
}
