package main

import "testing"

func TestConvertToFish(t *testing.T) {
	testCases := []struct {
		name       string
		bashInput  string
		fishOutput string
	}{
		{
			name:       "Conditional Operators",
			bashInput:  "echo 'Hello, World!' && echo 'Success' || echo 'Failure'",
			fishOutput: "echo 'Hello, World!' ; and echo 'Success' ; or echo 'Failure'",
		},
		{
			name:       "Command Substitution",
			bashInput:  "echo $(ls)",
			fishOutput: "echo (ls)",
		},
		{
			name:       "Backticks",
			bashInput:  "echo `ls`",
			fishOutput: "echo (ls)",
		},
		{
			name:       "Export",
			bashInput:  "export VAR=value",
			fishOutput: "set -x set VAR value",
		},
		{
			name:       "Variable Assignments",
			bashInput:  "VAR=value; echo $VAR",
			fishOutput: "set VAR value; echo $VAR",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fishOutput := convertToFish(tc.bashInput)
			if fishOutput != tc.fishOutput {
				t.Errorf("Expected: %s, Got: %s", tc.fishOutput, fishOutput)
			}
		})
	}
}
