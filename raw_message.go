package main

import (
	"encoding/json"

	"github.com/pkg/errors"
)

func getStringFromRawMessage(r json.RawMessage) (string, error) {
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

func getBoolFromRawMessage(r json.RawMessage) (bool, error) {
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

func getInt32FromRawMessage(r json.RawMessage) (int32, error) {
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

func getSliceFromRawMessage(r json.RawMessage) ([]json.RawMessage, error) {
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

func getMapFromRawMessage(r json.RawMessage) (map[string]json.RawMessage, error) {
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

func getSliceOfMapsFromRawMessage(r json.RawMessage) ([]map[string]json.RawMessage, error) {
	b, err := r.MarshalJSON()
	if err != nil {
		return nil, errors.Wrapf(err, "MarshalJSON(%s) failed: %s", r, err)
	}
	var sm []json.RawMessage
	if err = json.Unmarshal(b, &sm); err != nil {
		return nil, errors.Wrapf(err, "json.Unmarshal(%s failed: %s", b, err)
	}

	var result []map[string]json.RawMessage
	for _, m := range sm {
		c, err := m.MarshalJSON()
		if err != nil {
			return nil, errors.Wrapf(err, "MarshalJSON(%s) failed: %s", m, err)
		}
		var x map[string]json.RawMessage
		if err = json.Unmarshal(c, &x); err != nil {
			return nil, errors.Wrapf(err, "json.Unmarshal(%s failed: %s", c, err)
		}
		result = append(result, x)
	}

	return result, nil
}
