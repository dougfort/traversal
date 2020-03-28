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

func TestDictKey(t *testing.T) {
	var err error
	var buffer bytes.Buffer
	data := `{
		"name": "tagging",
		"category": "http"
	}`

	// valid key
	err = tr.Start([]byte(data)).DictKey("category").End(&buffer)
	if err != nil {
		t.Fatalf("error from valid JSON: %s", err)
	}

	if buffer.String() != "\"http\"" {
		t.Fatalf("invalid output: expected '\"http\"', found '%s'", buffer.String())
	}

	// invalid key
	err = tr.Start([]byte(data)).DictKey("foot").End(&buffer)
	if err == nil {
		t.Fatal("expecting error for invalid key")
	}

}

func TestListSingleton(t *testing.T) {
	var err error
	var buffer bytes.Buffer
	data := `[
		{"name": "tagging"}
	]`

	// valid key
	err = tr.Start([]byte(data)).ListSingleton().End(&buffer)
	if err != nil {
		t.Fatalf("error from valid JSON: %s", err)
	}

	if buffer.String() != "{\"name\": \"tagging\"}" {
		t.Fatalf("invalid output: expected '\"{\"name\": \"tagging\"}\"', found '%s'", buffer.String())
	}

	// empty list
	err = tr.Start([]byte("[]")).ListSingleton().End(&buffer)
	if err == nil {
		t.Fatal("expecting error for empty list")
	}

	// too big list
	bigData := `[
		{"name": "tagging"},
		{"category": "http"}
	]`
	err = tr.Start([]byte(bigData)).ListSingleton().End(&buffer)
	if err == nil {
		t.Fatal("expecting error too big list")
	}
}

func TestListPredicate(t *testing.T) {
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
	err = tr.Start([]byte(data)).ListPredicate(predicate).End(&buffer)
	if err != nil {
		t.Fatalf("error from valid JSON: %s", err)
	}

	if buffer.String() != "{\"key3\": \"value3\"}" {
		t.Fatalf("invalid output: expected '\"{\"key3\": \"value3\"}\"', found '%s'", buffer.String())
	}

	// empty list
	err = tr.Start([]byte("[]")).ListPredicate(predicate).End(&buffer)
	if err == nil {
		t.Fatal("expecting error for empty list")
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
		DictKey("cars").
		ListPredicate(predicate).
		DictKey("models").
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
		DictKey("configs").
		ListPredicate(func(r json.RawMessage) bool {
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
		DictKey("bootstrap").
		DictKey("static_resources").
		DictKey("listeners").
		ListSingleton().
		DictKey("filter_chains").
		ListSingleton().
		DictKey("filters").
		ListSingleton().
		DictKey("typed_config").
		DictKey("http_filters").
		ListPredicate(func(r json.RawMessage) bool {
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
		DictKey("typed_config").
		DictKey("value").
		End(os.Stdout)
	if err != nil {
		t.Fatalf("Traversal failed: %s", err)
	}
}
