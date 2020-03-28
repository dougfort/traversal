package main_test

import (
	"encoding/json"
	"fmt"
	"testing"

	tr "github.com/dougfort/traversal"
)

func TestGetStringFromRawMessage(t *testing.T) {
	var err error
	const testString = "test string"
	testJSON := fmt.Sprintf("\"%s\"", testString)

	var rawMessage json.RawMessage
	err = rawMessage.UnmarshalJSON([]byte(testJSON))
	if err != nil {
		t.Fatalf("rawMessage.UnmarshalJSON() failed: %s", err)
	}

	result, err := tr.GetStringFromRawMessage(rawMessage)
	if err != nil {
		t.Fatalf("GetStringFromRawMessage(%s) failed: %s", rawMessage, err)
	}
	if result != testString {
		t.Fatalf("expected '%s', got '%s", testString, result)
	}
}
