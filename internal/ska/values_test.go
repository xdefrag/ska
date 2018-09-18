package ska

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestParseValues(t *testing.T) {
	exs, err := ioutil.ReadDir(tdRaw)
	if err != nil {
		t.Fatal(err)
	}

	for _, ex := range exs {
		if !ex.IsDir() {
			continue
		}

		valuefile := tdRaw + ex.Name() + string(filepath.Separator) + "values.toml"

		got, err := ParseValues(valuefile)
		if err != nil {
			t.Fatal(err)
		}

		var want Values
		_, err = toml.DecodeFile(valuefile, &want)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}
}
