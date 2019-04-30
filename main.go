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
			must(err)

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

func tplPaths(ska, tpl string) (vp, tp string) {
	tplf := fmt.Sprintf("%s/%s", ska, tpl)

	return fmt.Sprintf("%s/values.toml", tplf), fmt.Sprintf("%s/templates", tplf)
}

func vals(path string) (map[string]interface{}, error) {
	var vals map[string]interface{}
	if _, err := toml.DecodeFile(path, &vals); err != nil {
		return nil, err
	}

	return vals, nil
}

func walk(in, out string, vals map[string]interface{}, f func(in, out string, vals map[string]interface{}) error) error {
	return filepath.Walk(in, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		saveto := out + string(filepath.Separator) + strings.Replace(path, in, "", -1)

		// if the filepath itself has templating, run it
		if strings.Contains(saveto, "{{") {
			saveto, err = genPath(saveto, vals)
			if err != nil {
				return err
			}
		}

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

// Generate a templated filename, rendering templated file/pathnames as needed
func genPath(path string, vals map[string]interface{}) (string, error) {
	t, err := template.New(path).Funcs(sprig.FuncMap()).Parse(path)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer([]byte(""))
	if err := t.Execute(buf, vals); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func mkdirr(path string) error {
	return os.MkdirAll(path, 0755)
}

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

func invokeEditor(ed, p string) error {
	cmd := exec.Command(ed, p)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
}
