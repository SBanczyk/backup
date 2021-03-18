package model

import (
	"encoding/json"
	"io/ioutil"
)

type StagingPath struct {
	Path   string
	Shadow bool
}

type Staging struct {
	StagingFiles   []StagingPath
	DestroyedFiles []string
}

func LoadStaging(path string) (backup *Staging, err error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s Staging
	err = json.Unmarshal(file, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func SaveStaging(path string, object *Staging) error {
	b, err := json.Marshal(object)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return err
	}
	return err
}
