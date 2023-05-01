package main

import (
	"net"
	"strings"
	"sync"
	"testing"
)

// test cases for randomString function
// テストを書く理由について: test を書くことで、コードの動作を確認できる。また、コードの変更によって、意図しない動作が起きないことを確認できる。
func TestRandomString(t *testing.T) {
	length := 10
	randomStr := randomString(length)

	if len(randomStr) != length {
		t.Errorf("Expected string length %d, got %d", length, len(randomStr))
	}

	for _, char := range randomStr {
		if !strings.ContainsRune(charset, char) {
			t.Errorf("Invalid character '%c' found in random string", char)
		}
	}
}

func TestResolveName(t *testing.T) {
	testCases := []struct {
		name        string
		expectError bool
	}{
		{"3-shake.com", false},
		{"invalid-domain-name-1234567890.com", true},
	}

	var wg sync.WaitGroup
	for _, tc := range testCases {
		wg.Add(1)

		go func(tc struct {
			name        string
			expectError bool
		}) {
			defer wg.Done()

			_, err := net.LookupHost(tc.name)

			if tc.expectError && err == nil {
				t.Errorf("Expected error for name %s, but got no error", tc.name)
			} else if !tc.expectError && err != nil {
				t.Errorf("Expected no error for name %s, but got error: %v", tc.name, err)
			}
		}(tc)
	}

	wg.Wait()
}
