package traversal

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// GetStringFromRawMessage returns a string from a JSON string
//
// given:
//
// 		\"value\" (with quotes)
//
// expecting:
//
// 		value
//
func GetStringFromRawMessage(r json.RawMessage) (string, error) {
	b, err := r.MarshalJSON()
	if err != nil {
		return "", errors.Wrapf(err, "MarshalJSON(%s) failed: %s", r, err)
	}
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		return "", errors.Wrapf(err, "json.Unmarshal(%s failed: %s", b, err)
	}

	return s, nil
}

// GetBoolFromRawMessage returns a boolean from a JSON string
//
// given:
//
// 		true (string)
//
// expecting:
//
// 		true (boolean)
//
func GetBoolFromRawMessage(r json.RawMessage) (bool, error) {
	b, err := r.MarshalJSON()
	if err != nil {
		return false, errors.Wrapf(err, "MarshalJSON(%s) failed: %s", r, err)
	}
	var i bool
	if err = json.Unmarshal(b, &i); err != nil {
		return false, errors.Wrapf(err, "json.Unmarshal(%s failed: %s", b, err)
	}

	return i, nil
}

// GetInt32FromRawMessage returns a number from a JSON string
//
// given:
//
// 		42 (string)
//
// expecting:
//
// 		42 (int32)
//
func GetInt32FromRawMessage(r json.RawMessage) (int32, error) {
	b, err := r.MarshalJSON()
	if err != nil {
		return -1, errors.Wrapf(err, "MarshalJSON(%s) failed: %s", r, err)
	}
	var i int32
	if err = json.Unmarshal(b, &i); err != nil {
		return -1, errors.Wrapf(err, "json.Unmarshal(%s failed: %s", b, err)
	}

	return i, nil
}

// GetSliceFromRawMessage returns a slice of json.RawMessage from a JSON RawMessage
//
// given:
//
// 		[
//   		"a": "a1",
//   		"b": "b1"
// 		]
//
// expecting:
//
// 		[]json.RawMessage{
//  		json.RawMessage{"a": "a1"},
//    		json.RawMessage{"b": "b1"}
// 		}
//
func GetSliceFromRawMessage(r json.RawMessage) ([]json.RawMessage, error) {
	b, err := r.MarshalJSON()
	if err != nil {
		return nil, errors.Wrapf(err, "MarshalJSON(%s) failed: %s", r, err)
	}
	var m []json.RawMessage
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, errors.Wrapf(err, "json.Unmarshal(%s failed: %s", b, err)
	}

	return m, nil
}

// GetMapFromRawMessage returns a map[string]json.RawMessage from a JSON RawMessage
//
// given:
//
// 		{
//   		"a": "a1",
//   		"b": "b1"
// 		}
//
// expecting:
//
// 		map[string]json.RawMessage{
//  		"a": json.RawMessage{"a1"},
//    		"b": json.RawMessage{"b1"}
// 		}
//
func GetMapFromRawMessage(r json.RawMessage) (map[string]json.RawMessage, error) {
	b, err := r.MarshalJSON()
	if err != nil {
		return nil, errors.Wrapf(err, "MarshalJSON(%s) failed: %s", r, err)
	}
	var m map[string]json.RawMessage
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, errors.Wrapf(err, "json.Unmarshal(%s failed: %s", b, err)
	}

	return m, nil
}
