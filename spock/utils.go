package spock

import (
	"io/ioutil"
	"path/filepath"
)

func combineConfigFiles(dir string) []byte {
	files, filesError := filepath.Glob(dir + "*.yml")
	if filesError != nil {
		return []byte{}
	}

	config := []byte{}
	for _, file := range files {
		contents, _ := ioutil.ReadFile(file)
		config = append(config, []byte("\n")...)
		config = append(config, contents...)
	}

	return config
}
