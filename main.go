package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/sprig"
)

func main() {
	fileVal := flag.String("values", "./values.toml", "Files with values for template")
	fileTpl := flag.String("template", "./template.tpl", "Template file or dir")
	fileOutput := flag.String("output", "./output", "Output")

	flag.Parse()

	in, out := *fileTpl, *fileOutput
	vals, err := vals(*fileVal)
	must(err)

	isdir, err := isDir(*fileTpl)
	must(err)

	if isdir {
		must(walk(in, out, vals, gen))
	} else {
		must(gen(in, out, vals))
	}
}

func vals(path string) (map[string]interface{}, error) {
	var vals map[string]interface{}
	if _, err := toml.DecodeFile(path, &vals); err != nil {
		return nil, err
	}

	return vals, nil
}

func isDir(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return f.IsDir(), nil
}

func walk(in, out string, vals map[string]interface{}, f func(in, out string, vals map[string]interface{}) error) error {
	return filepath.Walk(in, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		saveto := out + string(filepath.Separator) + strings.Replace(path, in, "", -1)

		if err := mkdirr(filepath.Dir(saveto)); err != nil {
			return err
		}

		return f(path, saveto, vals)
	})
}

func gen(in, out string, vals map[string]interface{}) error {
	t, err := template.New(filepath.Base(in)).Funcs(sprig.FuncMap()).ParseFiles(in)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer([]byte(""))
	if err := t.Execute(buf, vals); err != nil {
		return err
	}

	return ioutil.WriteFile(out, buf.Bytes(), 0644)
}

func mkdirr(path string) error {
	return os.MkdirAll(path, 0755)
}

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
}
