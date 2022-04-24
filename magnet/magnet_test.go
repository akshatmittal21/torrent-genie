package magnet

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/akshatmittal21/torrent-genie/constants"
	"github.com/akshatmittal21/torrent-genie/logger"
	"github.com/stretchr/testify/assert"
)

func TestGetLink(t *testing.T) {
	magnet := mockMagnet(t)
	link := magnet.GetLink("9128798362983", "test name")
	if link == "" {
		t.Error("link is empty")
	}
	if strings.Contains(link, "${INFO_HASH}") || strings.Contains(link, "${NAME}") || strings.Contains(link, "${TRACKER}") {
		t.Error("Link contains not expected values")
	}
}

func TestGetFile(t *testing.T) {
	magnet := mockMagnet(t)
	file := magnet.GetFile("9128798362983")

	if file == nil {
		t.Logf("file is nil")
	}
	filePath := constants.LogPath
	os.RemoveAll(path.Dir(path.Dir(filePath)))

}

func mockMagnet(t *testing.T) Magnet {
	filePath := constants.LogPath
	l, err := logger.Init(filePath, logger.DebugLevel)
	assert.NoError(t, err)
	return NewServer(l)
}
