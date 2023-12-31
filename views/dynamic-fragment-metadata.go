package views

import (
	"encoding/json"
	"os"
)

type templateMetafile struct {
	ResourceName string `json:"resourceName"`
	Hash         string `json:"hash"`
}

func loadFragmentMetafile(filepath string) (*templateMetafile, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var metafile templateMetafile
	err = json.NewDecoder(file).Decode(&metafile)
	if err != nil {
		return nil, err
	}

	return &metafile, nil
}
