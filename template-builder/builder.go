package templatebuilder

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/ncpa0/hardwire/configuration"
	"github.com/ncpa0/hardwire/utils"
)

func BuildPages(entrypoint string, outDir string, staticDir string, staticUrl string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if !path.IsAbs(entrypoint) {
		entrypoint = path.Join(wd, entrypoint)
	}
	if !path.IsAbs(outDir) {
		outDir = path.Join(wd, outDir)
	}
	if !path.IsAbs(staticDir) {
		staticDir = path.Join(wd, staticDir)
	}

	pagesDir := path.Dir(entrypoint)

	initProject(pagesDir)

	install := utils.Execute(("bun"), []string{
		"a",
		"hardwire-html-generator@0.0.1-beta.8", // Remember to update version after publish
	}, &utils.ExecuteOptions{
		Wd: pagesDir,
	})

	if install.Err != nil {
		return fmt.Errorf("error installing html builder package:\n%s %s", install.Stdout, install.Stderr)
	}

	installDev := utils.Execute(("bun"), []string{
		"a",
		"-D",
		"@types/bun",
	}, &utils.ExecuteOptions{
		Wd: pagesDir,
	})

	if installDev.Err != nil {
		return fmt.Errorf("error installing html builder package:\n%s %s", install.Stdout, install.Stderr)
	}

	builderInit := utils.Execute("bun", []string{
		"x",
		"hardwire-html-generator",
		"init",
		"--dir", pagesDir,
	}, &utils.ExecuteOptions{
		Wd: pagesDir,
	})

	if builderInit.Err != nil {
		return fmt.Errorf("error installing html builder package:\n%s %s", builderInit.Stdout, builderInit.Stderr)
	}

	if configuration.Current.DebugMode {
		fmt.Print("Building static HTML...\n")
	}

	result := utils.Execute("bun", []string{
		"x",
		"hardwire-html-generator",
		"build",
		"--src", entrypoint,
		"--outdir", outDir,
		"--staticdir", staticDir,
		"--staticurl", staticUrl,
	}, &utils.ExecuteOptions{
		Wd: pagesDir,
		Env: map[string]string{
			"NODE_ENV": "production",
		},
	})

	if result.Err != nil {
		return fmt.Errorf("error building pages:\n%s %s", result.Stdout, result.Stderr)
	}

	if configuration.Current.DebugMode {
		fmt.Printf("%s\n", result.Stdout)
	}

	return nil
}

type PackageJson struct {
	Name            string            `json:"name"`
	DevDependencies map[string]string `json:"devDependencies"`
	Main            string            `json:"main"`
}

func initProject(srcpath string) error {
	if _, err := os.Stat(srcpath); os.IsNotExist(err) {
		err := os.MkdirAll(srcpath, 0755)
		if err != nil {
			return err
		}
	}

	pkgJsonPath := path.Join(srcpath, "package.json")
	if _, err := os.Stat(pkgJsonPath); os.IsNotExist(err) {

		if configuration.Current.DebugMode {
			fmt.Print("Initializing the templates project\n")
		}

		pkgJson := PackageJson{
			Name: "page-templates",
			DevDependencies: map[string]string{
				"typescript": "latest",
			},
			Main: "index.tsx",
		}

		content, err := json.Marshal(pkgJson)

		if err != nil {
			return err
		}

		err = os.WriteFile(pkgJsonPath, content, 0644)
		if err != nil {
			return err
		}

		install := utils.Execute(("bun"), []string{
			"i",
		}, nil)

		if install.Err != nil {
			return fmt.Errorf("error installing html builder package:\n%s %s", install.Stdout, install.Stderr)
		}
	}
	return nil
}
