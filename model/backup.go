package model

import (
	"encoding/json"
	"io/ioutil"
)

type BackupPath struct {
	Path   string
	Shadow bool
}

type Backup struct {
	Files map[string][]BackupPath
	Version string
}

func LoadBackup(path string) (backup *Backup, err error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var b Backup
	err = json.Unmarshal(file, &b)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func SaveBackup(path string, object *Backup) error {
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
