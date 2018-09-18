package ska

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// GenerateTemplates from tmpldir to saveto with vv.
func GenerateTemplates(tmpldir, saveto string, vv Values) error {
	return filepath.Walk(tmpldir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		savetofile := saveto + string(filepath.Separator) + strings.Replace(path, tmpldir, "", -1)
		if err = ensureDirForFile(savetofile); err != nil {
			return errors.Wrapf(err, "Failed to create directory %s", savetofile)
		}

		t := template.New(filepath.Base(path))
		t, err = t.ParseFiles(path)
		if err != nil {
			return errors.Wrapf(err, "Failed to parse template %s", path)
		}

		buf := bytes.NewBuffer([]byte(""))
		if err = t.Execute(buf, vv); err != nil {
			return errors.Wrapf(err, "Failed to execute template %s", path)
		}

		if err = ioutil.WriteFile(savetofile, buf.Bytes(), 0755); err != nil {
			return errors.Wrapf(err, "Failed to write file %s", savetofile)
		}

		return nil
	})
}

func ensureDirForFile(path string) error {
	baseDir := filepath.Dir(path)
	_, err := os.Stat(baseDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return os.MkdirAll(baseDir, 0755)
}
