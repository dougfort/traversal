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

func TestGetBoolFromRawMessage(t *testing.T) {
	var err error
	testJSON := "true"

	var rawMessage json.RawMessage
	err = rawMessage.UnmarshalJSON([]byte(testJSON))
	if err != nil {
		t.Fatalf("rawMessage.UnmarshalJSON() failed: %s", err)
	}

	result, err := tr.GetBoolFromRawMessage(rawMessage)
	if err != nil {
		t.Fatalf("GetBoolFromRawMessage(%s) failed: %s", rawMessage, err)
	}
	if !result {
		t.Fatalf("expected 'true', got %t", result)
	}
}

func TestGetInt32FromRawMessage(t *testing.T) {
	var err error
	var testNumber int32 = 1423
	testJSON := fmt.Sprintf("%d", testNumber)

	var rawMessage json.RawMessage
	err = rawMessage.UnmarshalJSON([]byte(testJSON))
	if err != nil {
		t.Fatalf("rawMessage.UnmarshalJSON() failed: %s", err)
	}

	result, err := tr.GetInt32FromRawMessage(rawMessage)
	if err != nil {
		t.Fatalf("GetInt32FromRawMessage(%s) failed: %s", rawMessage, err)
	}
	if result != testNumber {
		t.Fatalf("expected %d, got %d", testNumber, result)
	}
}
