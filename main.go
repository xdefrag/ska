package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/sprig"
	"github.com/spf13/cobra"
)

func main() {
	var ska string
	var out string

	var cmd = &cobra.Command{
		Use:   "ska [template]",
		Short: "Render template",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			vp, tp := tplPaths(ska, args[0])

			vals, err := vals(vp)
			must(err)

			must(walk(tp, out, vals, gen))
		},
	}

	cmd.PersistentFlags().StringVarP(&ska, "templates", "t", "~/.local/share/ska", "Templates dir")
	cmd.PersistentFlags().StringVarP(&out, "output", "o", ".", "Output")

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
