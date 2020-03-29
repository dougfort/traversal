package traversal

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

// Traversal holds the state while travering JSON
type Traversal struct {
	err error
	msg json.RawMessage
}

// Start begins a Traversal by ijnitializing the internal state with JSON data.
func Start(data []byte) *Traversal {
	var t Traversal

	if json.Valid(data) {
		t.err = t.msg.UnmarshalJSON(data)
	} else {
		t.err = errors.Errorf("Invalid JSON")
	}

	return &t
}

// End terminates a Traversal
// If there is no error, End writes the internal state to the writer
// Note that a useful tool in debugging is to dump to os.Stdout
func (t *Traversal) End(w io.Writer) error {
	if t.err != nil {
		return t.err
	}

	data, err := t.msg.MarshalJSON()
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
	if t.err != nil {
		return t
	}

	m, err := GetMapFromRawMessage(t.msg)
	if err != nil {
		return &Traversal{
			err: errors.Wrap(err, "getMapFromRawMessage"),
			msg: nil,
		}
	}

	msg, ok := m[key]
	if !ok {
		return &Traversal{
			err: errors.Errorf("No entry for key '%s'", key),
			msg: nil,
		}
	}
	return &Traversal{err: nil, msg: msg}
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
	if t.err != nil {
		return t
	}

	s, err := GetSliceFromRawMessage(t.msg)
	if err != nil {
		return &Traversal{
			err: errors.Wrap(err, "getSliceFromRawMessage"),
			msg: nil,
		}
	}

	if len(s) != 1 {
		return &Traversal{
			err: errors.Errorf("Array has %d items", len(s)),
			msg: nil,
		}
	}

	return &Traversal{err: nil, msg: s[0]}
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
	if t.err != nil {
		return t
	}

	s, err := GetSliceFromRawMessage(t.msg)
	if err != nil {
		return &Traversal{
			err: errors.Wrap(err, "getSliceFromRawMessage"),
			msg: nil,
		}
	}

	for _, msg := range s {
		if p(msg) {
			return &Traversal{err: nil, msg: msg}
		}
	}

	return &Traversal{
		err: errors.Errorf("No Array item satisfies the predicate: %s", t.msg),
		msg: nil,
	}
}

// Selector selects whatever you want from the current traversal
//
// The go language does not allow you to add methods to an object outside of its package
// This gives you the ability to do just that
// We ask that you adhere to the spirit of composotion
// * No side effects
// * Pass on some proper subset of the incoming json.RawMessage, without making changes to it
//
// for example, you could duplicate ArrayPredicate
// given:
//
// [
//	{"key1": "value1"},
//	{"key2": "value2"},
//	{"key3": "value3"}
// ]
//
// with selector
//
//  func(r json.RawMessage) (json.RawMessage, error) {
//		s, err := tr.GetSliceFromRawMessage(r)
//		if err != nil {
//			return nil, err
//		}
//
//		for _, msg := range s {
//			m, err := tr.GetMapFromRawMessage(msg)
//			if err != nil {
//				return nil, err
//			}
//			v, ok := m["key3"]
//			if ok {
//				return v, nil
//			}
//		}
//
// 		if we make it here, we didn't find what we are looking for
//		return nil, fmt.Errorf("not found")
// }
//
// expecting
//
// "value3"
//
func (t *Traversal) Selector(s func(json.RawMessage) (json.RawMessage, error)) *Traversal {
	if t.err != nil {
		return t
	}

	var result Traversal
	result.msg, result.err = s(t.msg)

	return &result
}
