package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
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
	log.SetFlags(0)

	cmd := &cobra.Command{
		Use:   "ska [template]",
		Short: "Render template",
		Args:  cobra.ExactArgs(1),
		Run:   run,
	}

	setUpFlags(cmd)

	must(cmd.Execute())
}

// setUpFlags prepares command with flags.
func setUpFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP("templates", "t", templatePathDefault(), "templates dir")
	cmd.PersistentFlags().StringP("output", "o", ".", "output")
	cmd.PersistentFlags().StringP("editor", "e", os.Getenv("EDITOR"), "editor")
	cmd.PersistentFlags().BoolP("default-values", "d", false, "use default values")
}

// run runs ska scenario.
func run(cmd *cobra.Command, args []string) {
	var (
		ska           = cmd.Flag("templates").Value.String()
		out           = cmd.Flag("output").Value.String()
		editor        = cmd.Flag("editor").Value.String()
		defaultValues = cmd.Flag("default-values").Value.String()
	)

	valuePath, templatesPath := tplPaths(ska, args[0])

	out, err := filepath.Abs(out)
	if err != nil {
		must(err)
	}

	var values map[string]interface{}

	switch {
	case defaultValues == "true":
		values = readDefaultValues(valuePath)
	default:
		values = readValuesFromTempFile(valuePath, editor)
	}

	if err != nil && !os.IsNotExist(err) {
		must(err)
	}

	must(walk(templatesPath, out, values, gen))
}

// readDefaultValues reads values from default template values.
func readDefaultValues(valuePath string) map[string]interface{} {
	valuePath, _ = filepath.Abs(valuePath)

	values, err := vals(valuePath)
	if err != nil {
		log.Fatalf("Error while parsing default values: %v", err)
	}

	return values
}

// readValuesFromTempFile creates tempfile, starts editor, waits for stdout control
// and parse values.
func readValuesFromTempFile(valuePath, editor string) map[string]interface{} {
	tmp, err := tempfile(valuePath)
	if err != nil {
		log.Fatalf("Error while creating temp value file: %v", err)
	}

	s := bufio.NewScanner(os.Stdin)

	var values map[string]interface{}

	for {
		must(invokeEditor(editor, tmp))

		var err error

		values, err = vals(tmp)

		if !(err != nil) {
			break
		}

		fmt.Printf("Error while parsing file: %v\n", err)
		s.Scan()
	}

	if err := os.RemoveAll(tmp); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}

	return values
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
func walk(in, out string, vals map[string]interface{},
	f func(in, out string, vals map[string]interface{}) error) error {
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

		in, err = filepath.Abs(in)
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

	wd, _ := os.Getwd()

	rel, _ := filepath.Rel(wd, out)
	if rel == "" {
		rel = out
	}

	_, err = os.Stat(out)

	switch {
	case !(err != nil):
		log.Printf("\texists: %v", rel)
	case os.IsNotExist(err):
		if err := ioutil.WriteFile(out, buf.Bytes(), 0644); err != nil {
			return err
		}

		log.Printf("\tcreated: %v", rel)
	default:
	}

	return nil
}

// prepareFilepath generate filepath with vals values, removes ".ska" extension if any.
func prepareFilepath(path string, vals map[string]interface{}) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

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

	in, err := os.Open(pabs)
	if err != nil {
		return "", err
	}

	out, err := os.Create(tmp)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return "", err
	}

	return tmp, err
}

// invokeEditor invokes $EDITOR and pass stdin/stdout/stderr in it.
func invokeEditor(ed, p string) error {
	if ed == "" {
		log.Printf("WARNING: The $EDITOR environment variable has not been set, assuming /usr/bin/vim")

		ed = "/usr/bin/vim"
	}

	cmd := exec.Command(ed, p)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// templatePathDefault returns default ska templates folder.
func templatePathDefault() string {
	skadef, err := os.UserHomeDir()
	if err != nil {
		skadef = "/usr/local/share/ska"
	} else {
		skadef = fmt.Sprintf("%s/.local/share/ska", skadef)
	}

	return skadef
}

// must checks error and exit program if any.
func must(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Fprintf(os.Stderr, "%s:%d: %v\n", filepath.Base(file), line, err)
		os.Exit(-1)
	}
}
