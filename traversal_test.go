package main_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	tr "github.com/dougfort/traversal"
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

func TestArraySingleton(t *testing.T) {
	var err error
	var buffer bytes.Buffer
	data := `[
		{"name": "tagging"}
	]`

	// valid key
	err = tr.Start([]byte(data)).ArraySingleton().End(&buffer)
	if err != nil {
		t.Fatalf("error from valid JSON: %s", err)
	}

	if buffer.String() != "{\"name\": \"tagging\"}" {
		t.Fatalf("invalid output: expected '\"{\"name\": \"tagging\"}\"', found '%s'", buffer.String())
	}

	// empty Array
	err = tr.Start([]byte("[]")).ArraySingleton().End(&buffer)
	if err == nil {
		t.Fatal("expecting error for empty Array")
	}

	// too big Array
	bigData := `[
		{"name": "tagging"},
		{"category": "http"}
	]`
	err = tr.Start([]byte(bigData)).ArraySingleton().End(&buffer)
	if err == nil {
		t.Fatal("expecting error too big Array")
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
	err = tr.Start([]byte(data)).ArrayPredicate(predicate).End(&buffer)
	if err != nil {
		t.Fatalf("error from valid JSON: %s", err)
	}

	if buffer.String() != "{\"key3\": \"value3\"}" {
		t.Fatalf("invalid output: expected '\"{\"key3\": \"value3\"}\"', found '%s'", buffer.String())
	}

	// empty Array
	err = tr.Start([]byte("[]")).ArrayPredicate(predicate).End(&buffer)
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
		ArrayPredicate(predicate).
		ObjectKey("models").
		End(&buffer)
	if err != nil {
		t.Fatalf("error from valid JSON: %s", err)
	}

	t.Logf("%s", buffer.String())
}

func TestTraversal(t *testing.T) {
	const testFilePath = "configs.json"
	data, err := ioutil.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("ioutil.ReadFile(%s) failed: %s", testFilePath, err)
	}

	err = tr.Start(data).
		ObjectKey("configs").
		ArrayPredicate(func(r json.RawMessage) bool {
			m, err := tr.GetMapFromRawMessage(r)
			if err != nil {
				return false
			}
			n, err := tr.GetStringFromRawMessage(m["@type"])
			if err != nil {
				return false
			}
			return n == "type.googleapis.com/envoy.admin.v3.BootstrapConfigDump"
		}).
		ObjectKey("bootstrap").
		ObjectKey("static_resources").
		ObjectKey("listeners").
		ArraySingleton().
		ObjectKey("filter_chains").
		ArraySingleton().
		ObjectKey("filters").
		ArraySingleton().
		ObjectKey("typed_config").
		ObjectKey("http_filters").
		ArrayPredicate(func(r json.RawMessage) bool {
			m, err := tr.GetMapFromRawMessage(r)
			if err != nil {
				return false
			}
			n, err := tr.GetStringFromRawMessage(m["name"])
			if err != nil {
				return false
			}
			return n == "gm.metrics"
		}).
		ObjectKey("typed_config").
		ObjectKey("value").
		End(os.Stdout)
	if err != nil {
		t.Fatalf("Traversal failed: %s", err)
	}
}
