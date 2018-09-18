package ska

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestGenerateTemplates(t *testing.T) {
	exs, err := ioutil.ReadDir(tdRaw)
	if err != nil {
		t.Fatal(err)
	}

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
