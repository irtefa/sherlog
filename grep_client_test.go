// grep_client_test
package main

import (
	"strings"
	"testing"
)

const (
	TEST_SERVER_ADDR = "192.17.11.27:8008"
)

func Test_rewriteKeyAndValue(t *testing.T) {
	// Test message with ^ and $
	message := rewriteKeyAndValue("^well$", "^two$")
	if strings.EqualFold(message, "^well:two$") {
		t.Log("Regex Test 1 Passed")
	} else {
		t.Error("FAILED: Regex Test 1")
	}

	// Test message with some *
	message = rewriteKeyAndValue("^we*ll$", "^tw*o$")
	if strings.EqualFold(message, "^we*ll:tw*o$") {
		t.Log("Regex Test 2 Passed")
	} else {
		t.Error("FAILED: Regex Test 2")
	}

	// Test other messages with less ^ and $
	message = rewriteKeyAndValue("we*ll$", "^tw*o")
	if strings.EqualFold(message, "^[^:]*we*ll:tw*o[^:]*$") {
		t.Log("Regex Test 3 Passed")
	} else {
		t.Error("FAILED: Regex Test 3")
	}

	// Test message with no regex
	message = rewriteKeyAndValue("well", "two")
	if strings.EqualFold(message, "^[^:]*well[^:]*:[^:]*two[^:]*$") {
		t.Log("Regex Test 4 Passed")
	} else {
		t.Error("FAILED: Regex Test 4")
	}
}

func Test_writeToServer(t *testing.T) {
	c := make(chan string)

	//Rare
	realAnswer := "test.1\ni:only appear once\n"
	boolResult := writeToServerHelper("i", "only", realAnswer, c)
	if boolResult {
		t.Log("Server Communication Test 1 Passed")
	} else {
		t.Error("FAILED: Server Communication 1")
	}

	// -------- More cases --------
	//Sometimes
	realAnswer = "test.1\nthey:call me sometimes\nthey:call me sometimes\nthey:call me sometimes\nthey:call me sometimes\nthey:call me sometimes\n"
	boolResult = writeToServerHelper("they", "call", realAnswer, c)
	if boolResult {
		t.Log("Server Communication Test 2 Passed")
	} else {
		t.Error("FAILED: Server Communication 2")
	}

	//Frequent
	realAnswer = "test.1\nfrequent: line they call me\nfrequent: line they call me\nfrequent: line they call me\nfrequent: line they call me\nfrequent: line they call me\nfrequent: line they call me\nfrequent: line they call me\nfrequent: line they call me\nfrequent: line they call me\nfrequent: line they call me\nfrequent: line they call me\nfrequent: line they call me\nfrequent: line they call me\nfrequent: line they call me\nfrequent: line they call me\nfrequent: line they call me\nfrequent: line they call me\nfrequent: line they call me\nfrequent: line they call me\nfrequent: line they call me\n"
	boolResult = writeToServerHelper("frequent", "call", realAnswer, c)
	if boolResult {
		t.Log("Server Communication Test 3 Passed")
	} else {
		t.Error("FAILED: Server Communication 3")
	}

	//Regex
	realAnswer = "test.1\ni:only appear once\n"
	boolResult = writeToServerHelper("^.*", "^only appear once$", realAnswer, c)
	if boolResult {
		t.Log("Server Communication Test 1 Passed")
	} else {
		t.Error("FAILED: Server Communication 1")
	}
}

func writeToServerHelper(key string, value string, result string, c chan string) bool {
	message := "test" + rewriteKeyAndValue(key, value)
	go writeToServer(TEST_SERVER_ADDR, message, c)
	serverResult := <-c
	return strings.EqualFold(serverResult, result)
}
