//+build mage

// To install mage, run go get github.com/magefile/mage.

package main

import (
	"github.com/magefile/mage/sh"
)

func Build() error {
	return sh.Run("go", "build", "-o", "./dist/ska", "./...")
}

func Test() error {
	return sh.Run("go", "test", "-coverprofile", "coverage.txt", "-covermode", "atomic", "./...")
}

func TestBenchmark() error {
	return sh.Run("go", "test", "-bench", ".", "./...")
}

func TestRace() error {
	return sh.Run("go", "test", "-race", "./...")
}
