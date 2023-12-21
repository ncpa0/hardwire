package views

import (
	"encoding/json"
	"os"
)

type pageMetafile struct {
	IsDynamic bool `json:"isDynamic"`
	Resources [](struct {
		Key string `json:"key"`
		Res string `json:"res"`
	}) `json:"resources"`
	RedirectURL    string `json:"redirectUrl"`
	ShouldRedirect bool   `json:"shouldRedirect"`
}

func loadPageMetafile(filepath string) (*pageMetafile, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var metafile pageMetafile
	err = json.NewDecoder(file).Decode(&metafile)
	if err != nil {
		return nil, err
	}

	return &metafile, nil
}
