package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	must(os.RemoveAll("out"))

	os.Exit(m.Run())
}

func TestCompile(t *testing.T) {
	tests := map[string]struct {
		path string
	}{
		"a basic file with templated contents": {
			path: "dockerfile",
		},
		"templated directory and filenames": {
			path: "filenames",
		},
	}

	for msg, test := range tests {
		t.Run(msg, func(t *testing.T) {
			ska := "examples"
			tpl := test.path
			goldenDir := path.Join("testdata", test.path)
			outDir := path.Join("out", test.path)

			vp, tp := tplPaths(ska, tpl)
			vals, err := vals(vp)
			must(err)
			must(walk(tp, outDir, vals, gen))

			err = filepath.Walk(goldenDir, func(goldenPath string, goldenInfo os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				outPath := path.Join("out", strings.TrimPrefix(goldenPath, "testdata"))

				if goldenInfo.IsDir() {
					outInfo, err := os.Stat(outPath)
					if err != nil {
						t.Fatalf("Error for expected output directory %s: %v", outPath, err)
					}
					if !outInfo.IsDir() {
						t.Fatalf("Found file not directory for expected output directory %s", outPath)
					}
				} else {
					if !hasSameContents(t, outPath, goldenPath) {
						t.Fatalf("Compiled and golden files %s vs %s not the same", outPath, goldenPath)
					}
				}

				return nil
			})

			if err != nil {
				t.Fatal(err)
			}
			os.RemoveAll(outDir)

		})
	}
}

func hasSameContents(t *testing.T, f1, f2 string) bool {
	b1, err := ioutil.ReadFile(f1)
	if err != nil {
		t.Fatal(err)
	}

	b2, err := ioutil.ReadFile(f2)
	if err != nil {
		t.Fatal(err)
	}

	return bytes.Compare(b1, b2) == 0
}
