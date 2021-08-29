package logger

import (
	"os"
	"path"
	"testing"

	"github.com/akshatmittal21/torrent-genie/constants"
)

func TestInitLogger(t *testing.T) {

	filePath := constants.LogPath
	Error("writing log")

	if _, err := os.Stat(path.Dir(filePath)); os.IsNotExist(err) {
		t.Error("error creating log dir")
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("error creating file")
	}
	Info("writing log")
	if sugar == nil {
		t.Error("error creating logger")
	}
	os.RemoveAll(path.Dir(path.Dir(filePath)))
}
