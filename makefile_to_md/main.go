package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Step 1: Retrieve the Makefile path from command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path_to_Makefile>")
		return
	}
	makefilePath := os.Args[1]

	// Step 2: Open the Makefile
	file, err := os.Open(makefilePath)
	if err != nil {
		fmt.Println("Error opening Makefile:", err)
		return
	}
	defer file.Close()

	// Step 3: Initialize Markdown output
	var markdownContent strings.Builder
	markdownContent.WriteString("# 📋 Makefile Overview\n\n")

	// Step 4: Read the Makefile line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasSuffix(line, ":") {
			// Found a target
			targetName := strings.TrimSuffix(line, ":")
			markdownContent.WriteString(fmt.Sprintf("## 🎯 `%s` Target\n", targetName))
			markdownContent.WriteString(
				fmt.Sprintf("Run the following command:\n\n```bash\nmake %s\n```\n\n", targetName),
			)
		} else if strings.HasPrefix(line, "\t") {
			// Found a command
			markdownContent.WriteString(fmt.Sprintf("  - 🛠 Command: `%s`\n\n", strings.TrimPrefix(line, "\t")))
		} else if strings.HasPrefix(line, "#") {
			// Found a comment
			markdownContent.WriteString(fmt.Sprintf("- 💬 Comment: %s\n\n", strings.TrimPrefix(line, "#")))
		} else if strings.Contains(line, "=") {
			// Found a variable
			markdownContent.WriteString(fmt.Sprintf("- 🌍 Variable: `%s`\n\n", line))
		} else {
			// Other lines
			markdownContent.WriteString(fmt.Sprintf("- 📝 Other: `%s`\n\n", line))
		}
	}

	// Step 5: Check for errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading Makefile:", err)
		return
	}

	// Step 6: Output to standard output
	fmt.Print(markdownContent.String())
}
