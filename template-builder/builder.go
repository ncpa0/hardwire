package templatebuilder

import (
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

	if configuration.Current.DebugMode {
		fmt.Print("Building static HTML...\n")
	}

	pagesDir := path.Dir(entrypoint)

	install := utils.Execute(("bun"), []string{
		"a",
		"hardwire-html-generator@0.0.1-beta.6", // Remember to update version after publish
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
