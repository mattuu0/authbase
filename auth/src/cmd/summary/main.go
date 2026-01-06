package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	filePath := os.Getenv("TEST_SUMMARY_FILE")
	if filePath == "" {
		fmt.Println("TEST_SUMMARY_FILE environment variable not set.")
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening summary file: %v\n", err)
		return
	}
	defer file.Close()

	fmt.Println("\n" + `╔══════════════════════════════════════════════════════════════╗`)
	fmt.Println("`║                      TEST EXECUTION SUMMARY                  ║`")
	fmt.Println("`╠══════════════════════════════════════════════╦═══════════════╣`")
	fmt.Println("`║ Test Name                                    ║ Status        ║`")
	fmt.Println("`╠══════════════════════════════════════════════╬═══════════════╣`")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) == 2 {
			status := parts[0]
			name := parts[1]
			// Limit name length for table formatting
			if len(name) > 44 {
				name = name[:41] + "..."
			}
			fmt.Printf("║ %-44s ║ %-13s ║\n", name, status)
		}
	}

	fmt.Println("`╚══════════════════════════════════════════════╩═══════════════╝`")

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading summary file: %v\n", err)
	}
}
