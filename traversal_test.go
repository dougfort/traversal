package traversal_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	tr "github.com/dougfort/traversal"
	"github.com/pkg/errors"
)

func TestStartAndEnd(t *testing.T) {
	var err error
	var buffer bytes.Buffer

	// invalid data
	err = tr.Start([]byte("}")).End(&buffer)
	if err == nil {
		t.Fatal("expecting error from invalid JSON")
	}

	// valid data
	err = tr.Start([]byte("{}")).End(&buffer)
	if err != nil {
		t.Fatalf("error from valid JSON: %s", err)
	}

	if buffer.String() != "{}" {
		t.Fatalf("invalid output: expected '', found '%s'", buffer.String())
	}
}

func TestObjectKey(t *testing.T) {
	var err error
	var buffer bytes.Buffer
	data := `{
		"name": "tagging",
		"category": "http"
	}`

	// valid key
	err = tr.Start([]byte(data)).ObjectKey("category").End(&buffer)
	if err != nil {
		t.Fatalf("error from valid JSON: %s", err)
	}

	if buffer.String() != "\"http\"" {
		t.Fatalf("invalid output: expected '\"http\"', found '%s'", buffer.String())
	}

	// invalid key
	err = tr.Start([]byte(data)).ObjectKey("foot").End(&buffer)
	if err == nil {
		t.Fatal("expecting error for invalid key")
	}

}

func TestArrayPredicate(t *testing.T) {
	var err error
	var buffer bytes.Buffer
	data := `[
		{"key1": "value1"},
		{"key2": "value2"},
		{"key3": "value3"}
		]`
	predicate := func(r json.RawMessage) bool {
		m, err := tr.GetMapFromRawMessage(r)
		if err != nil {
			return false
		}
		n, err := tr.GetStringFromRawMessage(m["key3"])
		if err != nil {
			return false
		}
		return n == "value3"
	}

	// valid key
	err = tr.Start([]byte(data)).ArraySlice().Filter(predicate).End(&buffer)
	if err != nil {
		t.Fatalf("error from valid JSON: %s", err)
	}

	if buffer.String() != "{\"key3\": \"value3\"}" {
		t.Fatalf("invalid output: expected '\"{\"key3\": \"value3\"}\"', found '%s'", buffer.String())
	}

	// empty Array
	err = tr.Start([]byte("[]")).Filter(predicate).End(&buffer)
	if err == nil {
		t.Fatal("expecting error for empty Array")
	}
}

func TestArraySlice(t *testing.T) {
	var err error
	names := []string{"tagging", "tag"}
	data := fmt.Sprintf(`[
		{"name": "%s"},
		{"name": "%s"}
	]`, names[0], names[1])

	// valid key
	traversal := tr.Start([]byte(data)).ArraySlice()
	if traversal.Err != nil {
		t.Fatalf("error from valid JSON: %s", traversal.Err)
	}

	if len(traversal.Array) != 2 {
		t.Fatalf("Unexpected length: expected %d, dound %d", 2, len(traversal.Array))
	}

	// Test to see if we can unmarshal the bytes from ArraySlice() into a struct
	type testJSON struct {
		Name string `json:"name"`
	}
	testArray := make([]testJSON, len(traversal.Array))
	for i := 0; i < len(traversal.Array); i++ {
		err = json.Unmarshal(traversal.Array[i], &testArray[i])
		if err != nil {
			t.Fatal(errors.Wrap(err, "%d: failed to unmarshal msg from raw message"))
		}
		if testArray[i].Name != names[i] {
			t.Fatalf("name mismatch: expected '%s', found '%s'", names[i], testArray[i].Name)
		}
	}

}

func TestSelector(t *testing.T) {
	var err error
	var buffer bytes.Buffer
	data := `[
		{"key1": "value1"},
		{"key2": "value2"},
		{"key3": "value3"}
		]`

	selector := func(r []json.RawMessage) ([]json.RawMessage, error) {
		var array []json.RawMessage

		for _, msg := range r {
			m, err := tr.GetMapFromRawMessage(msg)
			if err != nil {
				return nil, err
			}
			value, ok := m["key3"]
			if ok {
				array = append(array, value)
			}
		}

		if len(array) == 0 {
			return nil, tr.ErrorEmptyArray
		}

		return array, nil
	}

	// valid key
	err = tr.Start([]byte(data)).ArraySlice().Selector(selector).End(&buffer)
	if err != nil {
		t.Fatalf("error from valid JSON: %s", err)
	}

	if buffer.String() != "\"value3\"" {
		t.Fatalf("invalid output: expected '\"value3\"', found '%s'", buffer.String())
	}

	// empty Array
	err = tr.Start([]byte("[]")).Selector(selector).End(&buffer)
	if err == nil {
		t.Fatal("expecting error for empty Array")
	}
}

func TestExample(t *testing.T) {
	var err error
	var buffer bytes.Buffer
	data := `{
		"name":"John",
		"age":30,
		"cars": [
		  { "name":"Ford", "models":[ "Fiesta", "Focus", "Mustang" ] },
		  { "name":"BMW", "models":[ "320", "X3", "X5" ] },
		  { "name":"Fiat", "models":[ "500", "Panda" ] }
		]
	   }`
	predicate := func(r json.RawMessage) bool {
		m, err := tr.GetMapFromRawMessage(r)
		if err != nil {
			return false
		}
		n, err := tr.GetStringFromRawMessage(m["name"])
		if err != nil {
			return false
		}
		return n == "BMW"
	}
	err = tr.Start([]byte(data)).
		ObjectKey("cars").
		ArraySlice().
		Filter(predicate).
		ObjectKey("models").
		End(&buffer)
	if err != nil {
		t.Fatalf("error from valid JSON: %s", err)
	}

	t.Logf("%s", buffer.String())
}
