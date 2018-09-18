package ska

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

var (
	// Editor that will be invoked with invokeEditor func.
	Editor = os.Getenv("EDITOR")
)

// Values type
type Values map[string]interface{}

// ParseValues from values file.
// Function will create temp file, open Editor, wait for edits and then parse values.
func ParseValues(valuesFilePath string) (Values, error) {
	temp, err := genTempFile(valuesFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create temp file")
	}

	defer os.RemoveAll(temp)

	vv := make(Values)

	for {
		invokeEditor(temp)
		_, err = toml.DecodeFile(temp, &vv)
		if err == nil {
			break
		}

		fmt.Println("Syntax error:")
		fmt.Println(err)
		fmt.Println()
		fmt.Println("Press ENTER to continue fix errors in values file or <C-c> to exit")
		bufio.NewScanner(os.Stdin).Scan()
	}

	return vv, nil
}

func genTempFile(path string) (string, error) {
	temp := ".temp-" + filepath.Base(path) + filepath.Ext(path)
	err := os.Link(path, temp)

	return temp, err
}

var invokeEditor = func(path string) error {
	cmd := exec.Command(Editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
