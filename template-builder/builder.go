package templatebuilder

import (
	"fmt"
	"os"
	"path"

	"github.com/ncpa0/hardwire/configuration"
	"github.com/ncpa0/hardwire/utils"
)

// //go:embed node_modules src bunfig.toml package.json
// var vfs embed.FS

// func extractFile(filename string, to string) error {
// 	outFile, err := os.OpenFile(to, os.O_CREATE|os.O_WRONLY, 0755)
// 	if err != nil {
// 		return err
// 	}
// 	defer outFile.Close()

// 	content, err := vfs.ReadFile(filename)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = outFile.Write(content)

// 	return err
// }

// func extractDir(dirpath string, to string) error {
// 	files, err := vfs.ReadDir(dirpath)
// 	if err != nil {
// 		return err
// 	}

// 	for _, file := range files {
// 		outfile := to + "/" + file.Name()
// 		if file.IsDir() {
// 			err = os.MkdirAll(outfile, 0755)
// 			if err != nil {
// 				return err
// 			}
// 			err = extractDir(dirpath+"/"+file.Name(), outfile)
// 			if err != nil {
// 				return err
// 			}
// 		} else {
// 			err = extractFile(dirpath+"/"+file.Name(), outfile)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// func extractVfs(entrypoint string, outDir string, staticDir string, staticUrl string) error {
// 	pagesDir := path.Dir(entrypoint)
// 	modulesDir := path.Join(pagesDir, "node_modules")
// 	binDir := path.Join(modulesDir, ".bin")
// 	packageDir := path.Join(modulesDir, "template-builder")
// 	subModulesDir := path.Join(packageDir, "node_modules")
// 	srcDir := path.Join(packageDir, "src")
// 	binFile := path.Join(srcDir, "index.tsx")
// 	bunfigFile := path.Join(packageDir, "bunfig.toml")
// 	pkgFile := path.Join(packageDir, "package.json")

// 	err = os.MkdirAll(outDir, 0755)
// 	if err != nil {
// 		return err
// 	}
// 	err = os.MkdirAll(staticDir, 0755)
// 	if err != nil {
// 		return err
// 	}

// 	err = os.MkdirAll(modulesDir, 0755)
// 	if err != nil {
// 		return err
// 	}
// 	err = os.MkdirAll(binDir, 0755)
// 	if err != nil {
// 		return err
// 	}
// 	err = os.MkdirAll(packageDir, 0755)
// 	if err != nil {
// 		return err
// 	}
// 	err = os.MkdirAll(subModulesDir, 0755)
// 	if err != nil {
// 		return err
// 	}
// 	err = os.MkdirAll(srcDir, 0755)
// 	if err != nil {
// 		return err
// 	}

// 	err = extractDir("node_modules", subModulesDir)
// 	if err != nil {
// 		return err
// 	}
// 	err = extractDir("src", srcDir)
// 	if err != nil {
// 		return err
// 	}
// 	err = extractFile("bunfig.toml", bunfigFile)
// 	if err != nil {
// 		return err
// 	}
// 	err = extractFile("package.json", pkgFile)
// 	if err != nil {
// 		return err
// 	}

// 	os.Symlink(binFile, path.Join(binDir, "template-builder"))

// 	return nil
// }

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

	// extractVfs(entrypoint, outDir, staticDir, staticUrl)
	pagesDir := path.Dir(entrypoint)

	install := utils.Execute(("bun"), []string{
		"a",
		"hardwire-html-generator@0.0.1-beta.5", // Remember to update version after publish
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
