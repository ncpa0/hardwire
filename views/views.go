package views

import (
	"fmt"
	"os"
	"path"
	"strings"

	config "github.com/ncpa0/hardwire/configuration"
	templatebuilder "github.com/ncpa0/hardwire/template-builder"
	"github.com/ncpa0/hardwire/utils"
	. "github.com/ncpa0cpl/ezs"
)

func IsTemplate(filepath string) bool {
	return strings.HasSuffix(filepath, ".template.html")
}

var pageViewRegistry = NewViewRegistry()
var dynamicFragmentViewRegistry = NewDynamicFragmentViewRegistry()

func LoadViews(wd string) error {
	htmlDir := config.Current.HtmlDir
	if !path.IsAbs(htmlDir) {
		htmlDir = path.Join(wd, htmlDir)
	}

	if !config.Current.NoBuild {
		if config.Current.CleanBuild {
			err := os.RemoveAll(htmlDir)
			if err != nil {
				return err
			}
		}

		err := templatebuilder.BuildPages(
			config.Current.Entrypoint,
			htmlDir,
			config.Current.StaticDir,
			config.Current.StaticURL,
		)

		if err != nil {
			return err
		}
	}

	if config.Current.DebugMode {
		fmt.Printf("Loading view from %s\n", htmlDir)
	}
	err := utils.Walk(htmlDir, func(root string, dirs []string, files []string) error {
		for _, file := range files {
			ext := path.Ext(file)

			if ext != ".html" {
				continue
			}

			fullPath := path.Join(root, file)
			relToView := fullPath[len(htmlDir):]

			if IsTemplate(relToView) {
				view, err := NewDynamicFragmentView(htmlDir, relToView)
				if err != nil {
					return err
				}

				if config.Current.DebugMode {
					fmt.Printf("Loading view from file %s\n", file)
					fmt.Printf("  ROOT: %s PATH: %s\n", htmlDir, relToView)
				}

				dynamicFragmentViewRegistry.Register(view)
			} else {
				view, err := NewPageView(htmlDir, relToView)
				if err != nil {
					return err
				}

				if config.Current.DebugMode {
					fmt.Printf("Loading view from file %s\n", file)
					fmt.Printf("  ROOT: %s PATH: %s\n", htmlDir, relToView)
				}

				pageViewRegistry.Register(view)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error loading views.")
		return err
	}

	err = loadIslands(htmlDir)

	if err != nil {
		fmt.Println("Error loading island views.")
		return err
	}

	return nil
}

func GetPageViewRegistry() *PageViewRegistry {
	return pageViewRegistry
}

func GetDynamicFragmentViewRegistry() *DynamicFragmentViewRegistry {
	return dynamicFragmentViewRegistry
}

func GetIslands() *Array[*Island] {
	return islandsList
}
