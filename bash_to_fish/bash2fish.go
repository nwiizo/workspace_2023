package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func replaceBackticks(input string) string {
	var output strings.Builder
	var backtickStart int
	backtickFlag := false

	for i, char := range input {
		if char == '`' {
			if backtickFlag {
				output.WriteString("(")
				output.WriteString(input[backtickStart+1 : i])
				output.WriteString(")")
				backtickFlag = false
			} else {
				backtickStart = i
				backtickFlag = true
			}
		} else {
			if !backtickFlag {
				output.WriteRune(char)
			}
		}
	}

	return output.String()
}

func replaceVariableAssignments(fishCommand string) string {
	varAssignmentRegex := regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)=([^ =]+)\b`)
	return varAssignmentRegex.ReplaceAllString(fishCommand, "set $1 $2")
}

func convertToFish(bashCommand string) string {
	fishCommand := bashCommand

	// Replace && with ; and
	fishCommand = strings.ReplaceAll(fishCommand, "&&", "; and")

	// Replace || with ; or
	fishCommand = strings.ReplaceAll(fishCommand, "||", "; or")

	// Replace $() with ()
	fishCommand = strings.ReplaceAll(fishCommand, "$(", "(")
	fishCommand = strings.ReplaceAll(fishCommand, ")", ")")

	// Replace `` with ()
	fishCommand = replaceBackticks(fishCommand)

	// Replace export with set -x
	fishCommand = strings.ReplaceAll(fishCommand, "export ", "set -x ")

	// Replace variable assignments VAR=value with set VAR value
	fishCommand = replaceVariableAssignments(fishCommand)

	return fishCommand
}

// Main function
// Reads the Bash one-liner from stdin and prints the converted Fish one-liner to stdout
func main() {
	fmt.Println("Enter the Bash one-liner:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	bashOneLiner := scanner.Text()

	fishOneLiner := convertToFish(bashOneLiner)
	fmt.Println("Fish one-liner:")
	fmt.Println(fishOneLiner)
}
