package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	must(os.RemoveAll("out"))

	os.Exit(m.Run())
}

// TODO make it auto.
func TestCompile(t *testing.T) {
	t.Run("dockerfile", func(t *testing.T) {
		ska := "examples"
		tpl := "dockerfile"
		out := "out/dockerfile"

		vp, tp := tplPaths(ska, tpl)

		vals, err := vals(vp)
		must(err)

		must(walk(tp, out, vals, gen))

		outfile := fmt.Sprintf("%s/Dockerfile", out)
		golden := "testdata/dockerfile/Dockerfile"

		if !isSame(t, outfile, golden) {
			t.Fatal("Compiled and golden files not the same")
		}
	})
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
