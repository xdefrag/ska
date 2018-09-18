package ska

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	clearTempOrPanic()
	disableInvokingEditor()

	c := m.Run()

	clearTempOrPanic()
	os.Exit(c)
}

func clearTempOrPanic() {
	if err := os.RemoveAll(tdTemp); err != nil {
		panic(err)
	}
}

func disableInvokingEditor() {
	invokeEditor = func(path string) error {
		return nil
	}

}
