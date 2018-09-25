package ska

import (
	"io/ioutil"
	"os"
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

	// Testdata examples tests
	for _, ex := range exs {
		if !ex.IsDir() {
			continue
		}

		t.Run(ex.Name(), func(t *testing.T) {
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
		})
	}

	// Error tests
	errCases := []struct {
		name string
		pre  func()
	}{
		{
			name: "genTempFile returns error",
			pre: func() {
				genTempFile = func(path string) (string, error) {
					return "", &os.LinkError{}
				}
			},
		},
		{
			name: "Toml decoder return error",
			pre: func() {
				decodeFile = func(p string, v interface{}) error {
					return &os.PathError{}
				}
			},
		},
	}

	for _, c := range errCases {
		t.Run(c.name, func(t *testing.T) {
			c.pre()

			_, err := ParseValues("")

			if err == nil {
				t.Errorf("Error expected")
			}
		})
	}
}
