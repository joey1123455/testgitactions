package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Execute git log command to get commit messages
	cmd := exec.Command("git", "log", "--pretty=format:%s", "HEAD^..HEAD")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing git log: %v\n", err)
		return
	}

	// Split the output by newline to get individual commit messages
	commitMessages := strings.Split(string(output), "\n")

	// Open commit.txt for writing
	file, err := os.Create("commit.txt")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	// Write each commit message to commit.txt
	for _, message := range commitMessages {
		_, err := file.WriteString(message + "\n")
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}
	}

	fmt.Println("Commit messages written to commit.txt")
}
