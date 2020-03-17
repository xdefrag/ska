package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestMain(m *testing.M) {
	must(os.RemoveAll("out"))

	st := m.Run()

	must(os.RemoveAll("out"))

	os.Exit(st)
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
			cmd := &cobra.Command{}

			setUpFlags(cmd)

			tmust(t, cmd.PersistentFlags().Set("templates", ska))
			tmust(t, cmd.PersistentFlags().Set("output", concpath(out, dir.Name())))
			tmust(t, cmd.PersistentFlags().Set("default-values", "true"))

			run(cmd, []string{dir.Name()})

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

		if !hasSameContents(t, path, golden) {
			t.Fatal("Compiled and golden files not the same")
		}

		return nil
	}
}

func hasSameContents(t *testing.T, f1, f2 string) bool {
	b1, err := ioutil.ReadFile(f1)
	if err != nil {
		t.Fatal(err)
	}

	// if file is empty = nothing need to generate.
	if len(b1) == 1 {
		_, err = os.Stat(f2)
		if os.IsNotExist(err) {
			return true
		}
	}

	b2, err := ioutil.ReadFile(f2)
	if err != nil {
		t.Fatal(err)
	}

	return bytes.Equal(b1, b2)
}

func concpath(p1, p2 string) string {
	return fmt.Sprintf("%s/%s", p1, p2)
}

func tmust(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
