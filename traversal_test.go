package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestTraversal(t *testing.T) {
	const testFilePath = "configs.json"
	data, err := ioutil.ReadFile(testFilePath)
	if err != nil {
		t.Fatalf("ioutil.ReadFile(%s) failed: %s", testFilePath, err)
	}

	err = Start(data).
		DictKey("configs").
		ListPredicate(func(r json.RawMessage) bool {
			m, err := getMapFromRawMessage(r)
			if err != nil {
				return false
			}
			n, err := getStringFromRawMessage(m["@type"])
			if err != nil {
				return false
			}
			return n == "type.googleapis.com/envoy.admin.v3.BootstrapConfigDump"
		}).
		DictKey("bootstrap").
		DictKey("static_resources").
		DictKey("listeners").
		ListIndex(0).
		DictKey("filter_chains").
		ListIndex(0).
		DictKey("filters").
		ListIndex(0).
		DictKey("typed_config").
		DictKey("http_filters").
		ListPredicate(func(r json.RawMessage) bool {
			m, err := getMapFromRawMessage(r)
			if err != nil {
				return false
			}
			n, err := getStringFromRawMessage(m["name"])
			if err != nil {
				return false
			}
			return n == "gm.metrics"
		}).
		JSON(os.Stdout)
	if err != nil {
		t.Fatalf("Traversal failed: %s", err)
	}
}
