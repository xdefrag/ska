package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	must(os.RemoveAll("out"))

	os.Exit(m.Run())
}

func TestExamples(t *testing.T) {
	const (
		ska      = "examples"
		out      = "out"
		testdata = "testdata"
	)

	dirs, err := ioutil.ReadDir(ska)
	if err != nil {
		t.Fatal(err)
	}

	for _, dir := range dirs {
		t.Run(dir.Name(), func(t *testing.T) {
			vp, tp := tplPaths(ska, dir.Name())

			vals, err := vals(vp)
			must(err)

			must(walk(tp, concpath(out, dir.Name()), vals, gen))

			if err := filepath.Walk(
				concpath(out, dir.Name()),
				compare(t, testdata, dir.Name()),
			); err != nil {
				t.Error(err)
			}
		})
	}
}

func compare(t *testing.T, tddir, tpldir string) func(string, os.FileInfo, error) error {
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		golden := strings.Replace(
			path,
			concpath("out", tpldir),
			concpath(tddir, tpldir),
			-1,
		)

		if !isSame(t, path, golden) {
			t.Fatal("Compiled and golden files not the same")
		}

		return nil
	}
}

func isSame(t *testing.T, f1, f2 string) bool {
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

func concpath(p1, p2 string) string {
	return fmt.Sprintf("%s/%s", p1, p2)
}
