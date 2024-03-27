package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
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
	// TODO: retrieve time from db
	cmd := exec.Command("git", "log")

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
		regex := regexp.MustCompile(`^[a-zA-Z0-9@#$%^&*()-_+=! ]+$`)

		if regex.MatchString(line) {
			gitLogs = append(gitLogs, line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error:", err)
	}

	log := parseLogs(gitLogs)
	// Print the git log messages
	// for _, _ = range logsList {
	// 	// println(log.id)
	// 	// println(log.author.username)
	// 	// println(log.author.email)
	// 	// println(log.date.String())
	// 	// println(log.message)
	// 	// println()
	// 	print(".")
	// }
	// println()

	content := log[0].date.String()

	// Write the string to a file named "output.txt"
	err = os.WriteFile(".github/workflows/lastOp.txt", []byte(content), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func parseLogs(logsList []string) []*gitLog {

	logs := make([]*gitLog, 0)
	//DONE: Retrieve from file
	packageVersion := readVersion()
	for _, logLine := range logsList {
		var newLog gitLog
		if strings.HasPrefix(logLine, "commit") {
			newLog = gitLog{
				id: strings.TrimSpace(logLine[len("commit "):]),
			}
			logs = append(logs, &newLog)

		} else {
			idx := len(logs) - 1

			switch {
			case strings.HasPrefix(logLine, "Author:"):
				pattern := regexp.MustCompile(`Author: (\w+) <([^>]+)>`)
				matches := pattern.FindStringSubmatch(logLine)
				if len(matches) == 3 {
					logs[idx].author.username = matches[1]
					logs[idx].author.email = matches[2]
				}

			case strings.HasPrefix(logLine, "Date:"):
				timeString := strings.TrimSpace(logLine[len("Date:"):])
				layout := "Mon Jan 02 15:04:05 2006 -0700"

				// Parse the time string
				t, err := time.Parse(layout, timeString)
				if err != nil {
					fmt.Println("Error parsing time:", err)
				}

				logs[idx].date = t.UTC()

			default:
				regex := regexp.MustCompile(`^[a-zA-Z0-9@#$%^&*()-_+=! ]+$`)

				// Check if the string matches the regular expression
				if regex.MatchString(logLine) {
					commitMessage := strings.TrimSpace(logLine)

					// DONE: update package version
					updatePackageVersion(strings.ToLower(commitMessage), packageVersion)
					logs[idx].message = commitMessage

				}
			}
		}
	}
	// DONE: update package json version
	updateVersionNo(packageVersion)
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
	// println(commit)

	switch {
	case matchBugFixes(commit):
		// println("bug fix")
		if version.patch < maxVersion {
			version.patch++
		} else {
			version.minor++
			version.patch = 0
		}
	case matchFeature(commit):
		// println("feature")
		// println(commit)
		if version.minor < maxVersion {
			version.minor++
		} else {
			version.major++
			version.minor = 0
		}
	case matchBreakingChange(commit):
		// println(commit)
		version.major++
	case matchChore(commit):
		// println("chore")
		if version.patch < maxVersion {
			version.patch++
		} else {
			version.minor++
			version.patch = 0
		}
	default:
		if version.patch < maxVersion {
			version.patch++
		} else {
			version.minor++
			version.patch = 0
		}
	}
}

func matchBugFixes(comment string) bool {
	pattern := `^(bug)(\s)?(fix)(es\s*)?`
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(comment)
}

func matchFeature(comment string) bool {
	pattern := `^(feat(ure)?|new(\s|-)?feat(ure)?)`
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(comment)
}

func matchBreakingChange(commit string) bool {
	pattern := `^(break(ing)?)(\s|-)?(change)(es\s*)?`
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(commit)
}

func matchChore(commit string) bool {
	pattern := `^(chore)(s)?`
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(commit)
}

func readVersion() *packageVersion {
	file, err := os.Open("package.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return nil
	}
	defer file.Close()

	// Declare a map to store the JSON data
	var jsonData map[string]interface{}

	// Decode the JSON data
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&jsonData)
	if err != nil {
		fmt.Println("Error decoding JSON data:", err)
		return nil
	}

	// Read the value of the "version" key
	version, ok := jsonData["version"].(string)
	if !ok {
		fmt.Println("Error: 'version' is not a string")
		return nil
	}

	// Split the string using "." as the delimiter
	parts := strings.Split(version, ".")

	// Check if the split produced three parts
	if len(parts) != 3 {
		fmt.Println("Error: The input string does not have three parts.")
		return nil
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		println("while parsing major version:", err)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		println("while parsing minor version:", err)
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		println("while parsing patch version:", err)
	}

	return &packageVersion{
		major: major,
		minor: minor,
		patch: patch,
	}
}

func updateVersionNo(version *packageVersion) {
	fileData, err := os.ReadFile("package.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	// Declare a map to store the JSON data
	var jsonData map[string]interface{}

	// Unmarshal the JSON data into the map
	err = json.Unmarshal(fileData, &jsonData)
	if err != nil {
		fmt.Println("Error unmarshaling JSON data:", err)
		return
	}

	// Update the value associated with the "version" key
	jsonData["version"] = fmt.Sprintf("%d.%d.%d", version.major, version.minor, version.patch)
	// jsonData["version"] = "2.1.0"

	// Marshal the updated data back to JSON
	updatedData, err := json.MarshalIndent(jsonData, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON data:", err)
		return
	}

	// Write the JSON data to the file
	err = os.WriteFile("package.json", updatedData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
		return
	}

	fmt.Println("JSON file updated successfully.")
}
