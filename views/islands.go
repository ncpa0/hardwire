package views

import (
	"encoding/json"
	"os"
	"path"
	"strings"

	"github.com/ncpa0/hardwire/utils"
	. "github.com/ncpa0cpl/convenient-structures"
)

type Island struct {
	ID         string
	FragmentID string
	Type       string
}

var islandsList = &Array[*Island]{}

func loadIslands(wd string) error {
	islandsDir := path.Join(wd, "__islands")

	err := utils.Walk(islandsDir, func(root string, dirs []string, files []string) error {
		for _, file := range files {
			if strings.HasSuffix(file, ".meta.json") {
				fullPath := path.Join(root, file)
				file, err := os.Open(fullPath)
				if err != nil {
					return err
				}
				defer file.Close()

				island := Island{}
				err = json.NewDecoder(file).Decode(&island)
				if err != nil {
					return err
				}
				validateIslandType(island.Type)
				islandsList.Push(&island)
			}
		}

		return nil
	})

	return err
}

func validateIslandType(itype string) {
	switch itype {
	case "basic":
	case "list":
	default:
		panic("Invalid island type")
	}
}
