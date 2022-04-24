package logger

import (
	"os"
	"path"
	"testing"

	"github.com/akshatmittal21/torrent-genie/constants"
	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {

	filePath := constants.LogPath
	l, err := Init(filePath, DebugLevel)
	assert.NoError(t, err)
	l.Error("writing log")

	if _, err := os.Stat(path.Dir(filePath)); os.IsNotExist(err) {
		t.Error("error creating log dir")
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("error creating file")
	}
	l.Info("writing log")
	os.RemoveAll(path.Dir(path.Dir(filePath)))
}
