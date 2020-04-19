package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type Script []map[string]json.RawMessage

func loadScript(path string) (Script, error) {
	var err error
	var scriptData []byte

	if scriptData, err = ioutil.ReadFile(path); err != nil {
		return nil, errors.Wrapf(err, "ioutil.ReadFile(%s)", path)
	}

	var script []map[string]json.RawMessage
	if err = json.Unmarshal(scriptData, &script); err != nil {
		return nil, errors.Wrapf(err, "json.Unmarshal(%s)", scriptData)
	}

	return script, nil
}
