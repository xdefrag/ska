package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/sprig"
	"github.com/spf13/cobra"
)

func main() {
	var ska string
	var out string
	var editor string

	var cmd = &cobra.Command{
		Use:   "ska [template]",
		Short: "Render template",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			vp, tp := tplPaths(ska, args[0])

			var vv map[string]interface{}

			tmp, err := tempfile(vp)

			// if there is no errors with tempfile = invoke editor to edit it.
			if !(err != nil) {
				s := bufio.NewScanner(os.Stdin)

				for {
					must(invokeEditor(editor, tmp))

					vv, err = vals(tmp)

					if !(err != nil) {
						break
					}

					fmt.Printf("Error while parsing file: %v\n", err)
					s.Scan()
				}

				if err := os.RemoveAll(tmp); err != nil {
					fmt.Fprintf(os.Stderr, "%v", err)
				}
			}

			if err != nil && !os.IsNotExist(err) {
				must(err)
			}

			must(walk(tp, out, vv, gen))
		},
	}

	skadef, err := os.UserHomeDir()
	if err != nil {
		skadef = "/usr/local/share/ska"
	} else {
		skadef = fmt.Sprintf("%s/.local/share/ska", skadef)
	}

	cmd.PersistentFlags().StringVarP(&ska, "templates", "t", skadef, "templates dir")
	cmd.PersistentFlags().StringVarP(&out, "output", "o", ".", "output")
	cmd.PersistentFlags().StringVarP(&editor, "editor", "e", os.Getenv("EDITOR"), "editor")

	must(cmd.Execute())
}

// tplPaths returns values.toml path (vp) and templates dir path (tp).
func tplPaths(ska, tpl string) (vp, tp string) {
	tplf := fmt.Sprintf("%s/%s", ska, tpl)

	return fmt.Sprintf("%s/values.toml", tplf), fmt.Sprintf("%s/templates", tplf)
}

// vals decodes path and return map of values with error.
func vals(path string) (map[string]interface{}, error) {
	var vals map[string]interface{}

	if _, err := toml.DecodeFile(path, &vals); err != nil {
		return nil, err
	}

	return vals, nil
}

// walk walks through in dirs, parses filenames witg vals, generates files with f functions and vals values.
func walk(in, out string, vals map[string]interface{}, f func(in, out string, vals map[string]interface{}) error) error {
	return filepath.Walk(in, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := prepareFilepath(path, vals)
		if err != nil {
			return err

		}

		saveto := out + string(filepath.Separator) + strings.Replace(file, in, "", -1)

		if err := mkdirr(filepath.Dir(saveto)); err != nil {
			return err
		}

		return f(path, saveto, vals)
	})
}

// gen generates files with in templates on out path with vals values.
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

// prepareFilepath generate filepath with vals values, removes ".ska" extention if any.
func prepareFilepath(path string, vals map[string]interface{}) (string, error) {
	// if the filepath itself has templating, run it
	if strings.Contains(path, "{{") {
		t, err := template.New(path).Funcs(sprig.FuncMap()).Parse(path)
		if err != nil {
			return "", err
		}

		buf := bytes.NewBuffer([]byte(""))
		if err := t.Execute(buf, vals); err != nil {
			return "", err
		}

		path = buf.String()
	}

	// if filepath has ".ska" ext, remove it.
	if path[len(path)-4:] == ".ska" {
		path = path[0 : len(path)-4]
	}

	return path, nil
}

// mkdirr recursivelly creates directories.
func mkdirr(path string) error {
	return os.MkdirAll(path, 0755)
}

// tempfile creates tempfiles in os.TempDir.
func tempfile(p string) (string, error) {
	tmp := fmt.Sprintf("%s/temp-%s", os.TempDir(), filepath.Base(p))
	pabs, err := filepath.Abs(p)

	if err != nil {
		return "", err
	}

	if err := os.Link(pabs, tmp); err != nil {
		return "", err
	}

	return tmp, err
}

// invokeEditor invokes $EDITOR and pass stdin/stdout/stderr in it.
func invokeEditor(ed, p string) error {
	cmd := exec.Command(ed, p)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// must checks error and exit program if any.
func must(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Fprintf(os.Stderr, "%s:%d: %v\n", filepath.Base(file), line, err)
		os.Exit(-1)
	}
}
