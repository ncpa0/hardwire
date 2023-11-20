package views

import (
	"encoding/json"
	"os"
)

type TemplateMetafile struct {
	ResourceName string `json:"resourceName"`
	Hash         string `json:"hash"`
}

func loadMetafile(filepath string) (*TemplateMetafile, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var metafile TemplateMetafile
	err = json.NewDecoder(file).Decode(&metafile)
	if err != nil {
		return nil, err
	}

	return &metafile, nil
}
