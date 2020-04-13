// Package traversal for traversing JSON text
package traversal

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

var (
	ErrorEmptyArray = errors.Errorf("No entries in Array")
)

// Traversal holds the state while travering JSON
type Traversal struct {

	// Err passes on an error
	// If Err is not nil, a componenet should just pass on the state
	Err error

	// Array represents the top level of JSON, which may be a single entry
	// It should never be empty
	Array []json.RawMessage
}

// Start begins a Traversal by initializing the internal state with JSON data.
func Start(data []byte) *Traversal {
	var t Traversal

	if json.Valid(data) {
		t.Array = make([]json.RawMessage, 1)
		t.Err = t.Array[0].UnmarshalJSON(data)
	} else {
		t.Err = errors.Errorf("Invalid JSON")
	}

	return &t
}

// End terminates a Traversal
// If there is no error, End writes Array[0] to the io.Writer
// Note that a useful tool in debugging is to dump to os.Stdout
func (t *Traversal) End(w io.Writer) error {
	if t.Err != nil {
		return t.Err
	}

	data, err := t.Array[0].MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "io.Write")
	}

	_, err = io.Copy(w, bytes.NewReader(data))
	if err != nil {
		return errors.Wrap(err, "o.Copy")
	}

	return nil
}

// ObjectKey replaces each element of Array with the contents of an Object at key
// This function assumes that each element of Array can be marshalled to an Object
// and that each Object has a value for key
//
// given:
//
// 		{
// 			"name": "tagging",
// 			"category": "http"
// 		}
// 		key = "category"
//
// expecting
//
// 		\"http\" (the output is JSON)
func (t *Traversal) ObjectKey(key string) *Traversal {
	if t.Err != nil {
		return t
	}

	for i := 0; i < len(t.Array); i++ {
		m, err := GetMapFromRawMessage(t.Array[i])
		if err != nil {
			return &Traversal{
				Err:   errors.Wrapf(err, "%d: GetMapFromRawMessage", i),
				Array: nil,
			}
		}

		var ok bool
		t.Array[i], ok = m[key]
		if !ok {
			return &Traversal{
				Err:   errors.Errorf("%d: No entry for key '%s'", i, key),
				Array: nil,
			}
		}
	}

	return t
}

// ArraySlice replaces the current Array with a slice unmarshalled from Array[0]
// It will fail if source Array does not exactly one element
// It will fail if the resulting Array does not have any elements
//
// given:
//
// 		[
//			{"key": "value"}
// 		]
//
// expecting
//
// 		[
//			{"key": "value"}
//		]
//
func (t *Traversal) ArraySlice() *Traversal {
	if t.Err != nil {
		return t
	}
	var err error

	if len(t.Array) != 1 {
		return &Traversal{
			Err:   errors.Wrapf(err, "Expecting Array size of 1; found %d", len(t.Array)),
			Array: nil,
		}
	}

	t.Array, err = GetSliceFromRawMessage(t.Array[0])
	if err != nil {
		return &Traversal{
			Err:   errors.Wrap(err, "GetSliceFromRawMessage"),
			Array: nil,
		}
	}
	if len(t.Array) == 0 {
		return &Traversal{
			Err:   ErrorEmptyArray,
			Array: nil,
		}
	}

	return t
}

// Filter replaces the current Array
//
// given:
//
// 		[
//			{"key1": "value1"},
//			{"key2": "value2"},
//			{"key3": "value3"}
// 		]
//
// with predicate
//
// 		func(r json.RawMessage) bool {
//			m, err := GetMapFromRawMessage(r)
//	    	if err != nil {
//          	return false
//	    	}
//	    	n, err := GetStringFromRawMessage(m["key3"])
//	    	if err != nil {
//		    	return false
//      	}
//      	return n == "value3"
// 		}
//
// expecting
//
// 		{"key3": "value3"}
//
func (t *Traversal) Filter(p func(json.RawMessage) bool) *Traversal {
	if t.Err != nil {
		return t
	}

	var array []json.RawMessage
	for _, entry := range t.Array {
		if p(entry) {
			array = append(array, entry)
		}
	}

	if len(array) == 0 {
		return &Traversal{
			Err:   ErrorEmptyArray,
			Array: nil,
		}
	}

	return &Traversal{Err: nil, Array: array}
}

// Selector selects whatever you want from the current traversal
//
// The go language does not allow you to add methods to an object outside of its package
// Selector gives you the ability to do just that
// We ask that you adhere to the spirit of composotion
// * No side effects
// * Pass on some proper subset of the incoming json.RawMessage, without making changes to it
//
// for example, you could duplicate ArrayPredicate
// given:
//
// 		[
//			{"key1": "value1"},
//			{"key2": "value2"},
//			{"key3": "value3"}
// 		]
//
// with selector
//
//  	func(r json.RawMessage) (json.RawMessage, error) {
//			s, err := tr.GetSliceFromRawMessage(r)
//			if err != nil {
//				return nil, err
//			}
//
//			for _, msg := range s {
//				m, err := tr.GetMapFromRawMessage(msg)
//				if err != nil {
//					return nil, err
//				}
//				v, ok := m["key3"]
//				if ok {
//					return v, nil
//				}
//			}
//
// 			if we make it here, we didn't find what we are looking for
//			return nil, fmt.Errorf("not found")
// 		}
//
// expecting
//
// 		"value3"
//
func (t *Traversal) Selector(s func([]json.RawMessage) ([]json.RawMessage, error)) *Traversal {
	if t.Err != nil {
		return t
	}

	var result Traversal
	result.Array, result.Err = s(t.Array)

	return &result
}
