package htmxframework

import (
	"fmt"
	"path"

	templatebuilder "github.com/ncpa0/htmx-framework/template-builder"
	"github.com/ncpa0/htmx-framework/utils"
	"github.com/ncpa0/htmx-framework/views"
)

var viewRegistry = views.NewViewRegistry()

func loadViews(wd string) error {
	viewsFullPath := conf.ViewsDir
	if !path.IsAbs(viewsFullPath) {
		viewsFullPath = path.Join(wd, viewsFullPath)
	}

	err := templatebuilder.BuildPages(
		conf.Entrypoint,
		viewsFullPath,
		conf.StaticDir,
		conf.StaticURL,
	)

	if err != nil {
		return err
	}

	if conf.DebugMode {
		fmt.Printf("Loading view from %s\n", viewsFullPath)
	}
	err = utils.Walk(viewsFullPath, func(root string, dirs []string, files []string) error {
		for _, file := range files {
			ext := path.Ext(file)

			if ext != ".html" {
				continue
			}

			fullPath := path.Join(root, file)
			relToView := fullPath[len(viewsFullPath):]
			view, err := views.NewView(viewsFullPath, relToView)
			if err != nil {
				return err
			}

			if conf.DebugMode {
				fmt.Printf("Loading view from file %s\n", file)
				fmt.Printf("  ROOT: %s PATH: %s\n", viewsFullPath, relToView)
			}

			viewRegistry.Register(view)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error loading views.")
		return err
	}

	return nil
}
