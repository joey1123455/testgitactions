package main

import (
	"fmt"
	"os"
)

func main() {
	// Open the file in append mode. If the file doesn't exist, it will be created.
	file, err := os.OpenFile("commits.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Write "Hello, World!" to the file.
	if _, err := file.WriteString("Hello, World!\n"); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Successfully wrote to commits.txt")
}
