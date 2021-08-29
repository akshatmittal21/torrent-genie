package magnet

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/akshatmittal21/torrent-genie/constants"
)

func TestGetLink(t *testing.T) {
	link := GetLink("9128798362983", "test name")
	if link == "" {
		t.Error("link is empty")
	}
	if strings.Contains(link, "${INFO_HASH}") || strings.Contains(link, "${NAME}") || strings.Contains(link, "${TRACKER}") {
		t.Error("Link contains not expected values")
	}
}

func TestGetFile(t *testing.T) {
	file := GetFile("9128798362983")

	if file == nil {
		t.Logf("file is nil")
	}
	filePath := constants.LogPath
	os.RemoveAll(path.Dir(path.Dir(filePath)))

}
