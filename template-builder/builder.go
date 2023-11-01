package templatebuilder

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"os/exec"
)

//go:embed node_modules/jsxte node_modules/minimist node_modules/prettier src bunfig.toml
var vfs embed.FS

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

func BuildPages(pagesDir string, outDir string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	tmpDir := wd + "/.tmp"
	nodemodulesDir := tmpDir + "/node_modules"
	binFile := tmpDir + "/index.tsx"
	bunConfigFile := tmpDir + "/bunfig.toml"

	err = os.MkdirAll(tmpDir, 0755)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

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

	cmd := exec.Command("bun", "--config="+bunConfigFile, binFile, "build", "--src", pagesDir, "--outdir", outDir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	cmd.Env = append(cmd.Environ(), "NODE_ENV=production")
	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("error building pages: %s %s", stdout.String(), stderr.String())
	}

	return nil
}
