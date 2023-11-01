package templatebuilder

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path"
)

//go:embed src bunfig.toml node_modules/jsxte node_modules/minimist node_modules/prettier node_modules/csso  node_modules/css-tree node_modules/mdn-data node_modules/source-map-js
var vfs embed.FS
var DebugMode = false

func extractFile(filename string, to string) error {
	outFile, err := os.OpenFile(to, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer outFile.Close()

	content, err := vfs.ReadFile(filename)
	if err != nil {
		return err
	}

	_, err = outFile.Write(content)

	return err
}

func extractDir(dirpath string, to string) error {
	files, err := vfs.ReadDir(dirpath)
	if err != nil {
		return err
	}

	for _, file := range files {
		outfile := to + "/" + file.Name()
		if file.IsDir() {
			err = os.MkdirAll(outfile, 0755)
			if err != nil {
				return err
			}
			err = extractDir(dirpath+"/"+file.Name(), outfile)
			if err != nil {
				return err
			}
		} else {
			err = extractFile(dirpath+"/"+file.Name(), outfile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func BuildPages(pagesDir string, outDir string, staticDir string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if !path.IsAbs(pagesDir) {
		pagesDir = path.Join(wd, pagesDir)
	}
	if !path.IsAbs(outDir) {
		outDir = path.Join(wd, outDir)
	}
	if !path.IsAbs(staticDir) {
		staticDir = path.Join(wd, staticDir)
	}

	tmpDir := path.Join(wd, ".tmp")
	nodemodulesDir := path.Join(tmpDir, "node_modules")
	binFile := path.Join(tmpDir, "index.tsx")
	bunConfigFile := path.Join(tmpDir, "bunfig.toml")

	err = os.MkdirAll(tmpDir, 0755)
	if err != nil {
		return err
	}
	if !DebugMode {
		defer os.RemoveAll(tmpDir)
	}

	err = extractDir("node_modules", nodemodulesDir)
	if err != nil {
		return err
	}
	err = extractDir("src", tmpDir)
	if err != nil {
		return err
	}
	err = extractFile("bunfig.toml", bunConfigFile)
	if err != nil {
		return err
	}

	cmd := exec.Command("bun", "--config="+bunConfigFile, binFile, "build", "--src", pagesDir, "--outdir", outDir, "--staticdir", staticDir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmd.Dir = tmpDir
	cmd.Env = append(cmd.Environ(), "NODE_ENV=production")
	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("error building pages: %s %s", stdout.String(), stderr.String())
	}

	return nil
}
