package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	tr "github.com/dougfort/traversal"
)

func isObject(b byte) bool {
	return b == '{'
}

func isArray(b byte) bool {
	return b == '['
}

func main() {
	var err error
	var data []byte
	var script Script

	if len(os.Args) < 3 {
		log.Fatalf("You must specify path to data file and path to script")
	}

	path := os.Args[1]
	if data, err = ioutil.ReadFile(path); err != nil {
		log.Fatalf("ioutil.ReadFile(%s) failed: %s", path, err)
	}

	if !json.Valid(data) {
		log.Fatalf("invalid JSON: %s", path)
	}

	if script, err = loadScript(os.Args[2]); err != nil {
		log.Fatalf("loadScript(%s) failed: %s", os.Args[2], err)
	}

	var traversal *tr.Traversal
	for i, rawMap := range script {
		name, err := tr.GetStringFromRawMessage(rawMap["name"])
		if err != nil {
			log.Fatalf("tr.GetStringFromRawMessage(%s) failed: %s", rawMap["name"], err)
		}
		switch name {
		case "start":
			log.Printf("Start")
			traversal = tr.Start(data)
			if traversal.Err != nil {
				log.Fatalf("Traversal Err = %s", traversal.Err)
			}
			log.Printf("#%d: Start: %d", i+1, len(traversal.Array))
		case "object-key":
			key, err := tr.GetStringFromRawMessage(rawMap["key"])
			if err != nil {
				log.Fatalf("tr.GetStringFromRawMessage(%s) failed: %s", rawMap["key"], err)
			}
			traversal = traversal.ObjectKey(key)
			if traversal.Err != nil {
				log.Fatalf("Traversal Err = %s", traversal.Err)
			}
			log.Printf("#%d: ObjectKey(%s): %d", i+1, key, len(traversal.Array))
		case "array-slice":
			traversal = traversal.ArraySlice()
			if traversal.Err != nil {
				log.Fatalf("Traversal Err = %s", traversal.Err)
			}
			log.Printf("#%d: ArraySlice:, %d", i+1, len(traversal.Array))
		case "inspect":
			key, err := tr.GetStringFromRawMessage(rawMap["key"])
			if err != nil {
				log.Fatalf("tr.GetStringFromRawMessage(%s) failed: %s", rawMap["key"], err)
			}
			traversal.Inspect(func(r json.RawMessage) {
				m, err := tr.GetMapFromRawMessage(r)
				if err != nil {
					log.Printf("ERROR: %s", err)
				} else {
					value, err := tr.GetStringFromRawMessage(m[key])
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Println(value)
					}
				}
			})

			log.Printf("#%d: Inspect:, %d", i+1, len(traversal.Array))
		case "filter":
			key, err := tr.GetStringFromRawMessage(rawMap["key"])
			if err != nil {
				log.Fatalf("tr.GetStringFromRawMessage(%s) failed: %s", rawMap["key"], err)
			}
			value, err := tr.GetStringFromRawMessage(rawMap["value"])
			if err != nil {
				log.Fatalf("tr.GetStringFromRawMessage(%s) failed: %s", rawMap["value"], err)
			}

			traversal = traversal.Filter(func(r json.RawMessage) bool {
				m, err := tr.GetMapFromRawMessage(r)
				if err != nil {
					return false
				}

				v, err := tr.GetStringFromRawMessage(m[key])
				if err != nil {
					return false
				}

				return v == value
			})
			if traversal.Err != nil {
				log.Fatalf("Traversal Err = %s", traversal.Err)
			}

			log.Printf("#%d: Filter:, %d", i+1, len(traversal.Array))
		default:
			log.Printf("#%d: Unknown name: '%s'", i+1, name)
		}
	}

}
