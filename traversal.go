package main

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

// Traversal holds the state while travering JSON
type Traversal struct {
	Err error
	Msg json.RawMessage
}

// Start begins a Traversal by ijnitializing the internal state with JSON data.
func Start(data []byte) *Traversal {
	var t Traversal

	if json.Valid(data) {
		t.Err = t.Msg.UnmarshalJSON(data)
	} else {
		t.Err = errors.Errorf("Invalid JSON")
	}

	return &t
}

// End terminates a Traversal
// If there is no error, End writes the internal state to the writer
// Note that a useful tool in debugging is to dump to os.Stdout
func (t *Traversal) End(w io.Writer) error {
	if t.Err != nil {
		return t.Err
	}

	data, err := t.Msg.MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "io.Write")
	}

	_, err = io.Copy(w, bytes.NewReader(data))
	if err != nil {
		return errors.Wrap(err, "o.Copy")
	}

	return nil
}

// ObjectKey selects a key from a JSON object (go map)
//
// given:
//
// {
// 	"name": "tagging",
// 	"category": "http"
// }
// key = "category"
//
// expecting
//
// \"http\" (the output is JSON)
func (t *Traversal) ObjectKey(key string) *Traversal {
	if t.Err != nil {
		return t
	}

	m, err := GetMapFromRawMessage(t.Msg)
	if err != nil {
		return &Traversal{
			Err: errors.Wrap(err, "getMapFromRawMessage"),
			Msg: nil,
		}
	}

	msg, ok := m[key]
	if !ok {
		return &Traversal{
			Err: errors.Errorf("No entry for key '%s'", key),
			Msg: nil,
		}
	}
	return &Traversal{Err: nil, Msg: msg}
}

// ArraySingleton selects the only entry from an Array.
// It will fail if the Array does not have exactly one item.
//
// given:
//
// [
//	{"key": "value"}
// ]
//
// expecting
//
// {"key": "value"}
//
func (t *Traversal) ArraySingleton() *Traversal {
	if t.Err != nil {
		return t
	}

	s, err := GetSliceFromRawMessage(t.Msg)
	if err != nil {
		return &Traversal{
			Err: errors.Wrap(err, "getSliceFromRawMessage"),
			Msg: nil,
		}
	}

	if len(s) != 1 {
		return &Traversal{
			Err: errors.Errorf("Array has %d items", len(s)),
			Msg: nil,
		}
	}

	return &Traversal{Err: nil, Msg: s[0]}
}

// ArrayPredicate selects an entry from an Array based on a predicate
//
// given:
//
// [
//	{"key1": "value1"},
//	{"key2": "value2"},
//	{"key3": "value3"}
// ]
//
// with predicate
//
// func(r json.RawMessage) bool {
//		m, err := GetMapFromRawMessage(r)
//	    if err != nil {
//          return false
//	    }
//	    n, err := GetStringFromRawMessage(m["key3"])
//	    if err != nil {
//		    return false
//      }
//      return n == "value3"
// }
//
// expecting
//
// {"key3": "value3"}
//
func (t *Traversal) ArrayPredicate(p func(json.RawMessage) bool) *Traversal {
	if t.Err != nil {
		return t
	}

	s, err := GetSliceFromRawMessage(t.Msg)
	if err != nil {
		return &Traversal{
			Err: errors.Wrap(err, "getSliceFromRawMessage"),
			Msg: nil,
		}
	}

	for _, msg := range s {
		if p(msg) {
			return &Traversal{Err: nil, Msg: msg}
		}
	}

	return &Traversal{
		Err: errors.Errorf("No Array item satisfies the predicate: %s", t.Msg),
		Msg: nil,
	}
}
