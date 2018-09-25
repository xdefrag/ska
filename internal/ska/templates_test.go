package ska

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestGenerateTemplates(t *testing.T) {
	exs, err := ioutil.ReadDir(tdRaw)
	if err != nil {
		t.Fatal(err)
	}

	// Testdata examples tests
	for _, ex := range exs {

		valuefile := tdRaw + ex.Name() + string(filepath.Separator) + "values.toml"

		var vv Values
		_, err = toml.DecodeFile(valuefile, &vv)
		if err != nil {
			t.Fatal(err)
		}

		if err = GenerateTemplates(tdRaw+ex.Name()+string(filepath.Separator)+"templates/",
			tdTemp+ex.Name(),
			vv); err != nil {
			t.Fatal(err)
		}

		check(t, ex.Name())
	}

	// Error tests
	errCases := []struct {
		name string
		pre  func()
	}{
		{
			name: "Template dir doesn't exists error",
			pre: func() {
				restore := stat

				stat = func(path string) error {
					stat = restore

					return os.ErrNotExist
				}
			},
		},
		{
			name: "Create dir error",
			pre: func() {
				restore := ensureDirForFile

				ensureDirForFile = func(path string) error {
					ensureDirForFile = restore

					return os.ErrNotExist
				}
			},
		},
		{
			name: "Template execute error",
			pre: func() {
				restore := templateExecute

				templateExecute = func(path string, vv Values) (*bytes.Buffer, error) {
					templateExecute = restore

					return nil, &template.Error{}
				}
			},
		},
		{
			name: "Write file error",
			pre: func() {
				restore := writeFile

				writeFile = func(path string, data []byte) error {
					writeFile = restore

					return os.ErrExist
				}
			},
		},
	}

	for _, c := range errCases {
		t.Run(c.name, func(t *testing.T) {
			c.pre()

			if err = GenerateTemplates(tdRaw, tdTemp, Values{}); err == nil {
				t.Error("Error expected")
			}
		})
	}
}

func check(t *testing.T, path string) {
	files, err := ioutil.ReadDir(tdCompiled + path)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			check(t, path+string(filepath.Separator)+file.Name())
			continue
		}

		want, err := ioutil.ReadFile(tdCompiled + path + string(filepath.Separator) + file.Name())
		if err != nil {
			t.Fatal(err)
		}

		got, err := ioutil.ReadFile(tdTemp + path + string(filepath.Separator) + file.Name())
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}
}
