package main

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

// Traversal marks progress in traversing a JSON graph
type Traversal struct {
	err error
	msg json.RawMessage
}

// Start a Traversal with JSON bytes
func Start(data []byte) *Traversal {
	var t Traversal
	t.err = t.msg.UnmarshalJSON(data)
	return &t
}

// JSON copies the interal state to a writer
// Note that a useful too in debugging is to dump to os.Stdout
func (t *Traversal) JSON(w io.Writer) error {
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

// DictKey selects a key from a map
func (t *Traversal) DictKey(key string) *Traversal {
	if t.err != nil {
		return t
	}

	m, err := getMapFromRawMessage(t.msg)
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

// ListIndex selects an entry from a list
func (t *Traversal) ListIndex(index int) *Traversal {
	if t.err != nil {
		return t
	}

	s, err := getSliceFromRawMessage(t.msg)
	if err != nil {
		return &Traversal{
			err: errors.Wrap(err, "getSliceFromRawMessage"),
			msg: nil,
		}
	}

	if index >= len(s) {
		return &Traversal{
			err: errors.Errorf("Invalid index %d of %d", index, len(s)),
			msg: nil,
		}
	}

	return &Traversal{err: nil, msg: s[index]}
}

// ListPredicate selects an entry from a list based on a predicate
func (t *Traversal) ListPredicate(p func(json.RawMessage) bool) *Traversal {
	if t.err != nil {
		return t
	}

	s, err := getSliceFromRawMessage(t.msg)
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
		err: errors.Errorf("No list item satisfies the predicate: %s", t.msg),
		msg: nil,
	}
}
