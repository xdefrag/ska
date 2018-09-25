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

		buf, err := templateExecute(path, vv)
		if err != nil {
			return errors.Wrapf(err, "Failed to execute template %s", path)
		}

		if err = writeFile(savetofile, buf.Bytes()); err != nil {
			return errors.Wrapf(err, "Failed to write file %s", savetofile)
		}

		return nil
	})
}

var ensureDirForFile = func(path string) error {
	baseDir := filepath.Dir(path)
	_, err := os.Stat(baseDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return os.MkdirAll(baseDir, 0755)
}

var templateExecute = func(path string, vv Values) (buf *bytes.Buffer, err error) {
	buf = bytes.NewBuffer([]byte(""))

	t, err := template.New(filepath.Base(path)).ParseFiles(path)
	if err != nil {
		return
	}

	err = t.Execute(buf, vv)

	return
}

var writeFile = func(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0755)
}
