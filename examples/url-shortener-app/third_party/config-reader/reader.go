package config_reader

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Read(filename string) (settings map[string]any, err error) {
	filename = filepath.Clean(filename)
	// We get the full file name for a nice error in case of a file reading error.
	if filename, err = filepath.Abs(filename); err != nil {
		return nil, err
	}
	if data, err := os.ReadFile(filename); err != nil {
		return nil, err
	} else if err = json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}
	return settings, err
}
