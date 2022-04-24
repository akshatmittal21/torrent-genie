package torrent

import (
	"os"
	"path"
	"testing"

	"github.com/akshatmittal21/torrent-genie/constants"
	"github.com/akshatmittal21/torrent-genie/dto"
	"github.com/akshatmittal21/torrent-genie/logger"
	"github.com/akshatmittal21/torrent-genie/torrent/piratebay"
	"github.com/stretchr/testify/assert"
)

func TestGetTorrent(t *testing.T) {
	filePath := constants.LogPath
	l, err := logger.Init(filePath, logger.DebugLevel)
	assert.NoError(t, err)

	torr := NewServer(l)
	torrents, err := torr.GetTorrents("abcd")
	assert.NoError(t, err)
	assert.NotZero(t, len(torrents))

	torrents, err = torr.GetTorrents("oquu2gdbi")
	assert.NoError(t, err)
	assert.NotZero(t, len(torrents))
	assert.Equal(t, torrents[0].ID, "0")

	torrents, err = torr.GetTorrents("")
	assert.NoError(t, err)
	assert.NotZero(t, len(torrents))
	assert.Equal(t, torrents[0].ID, "0")
}

type recentTorrent struct {
	arg1     constants.Command
	expected []piratebay.Torrent
}

var recentTorrents = []recentTorrent{
	{constants.RECENT, make([]piratebay.Torrent, constants.RecentCount)},
	{constants.TOPAUDIO, make([]piratebay.Torrent, constants.RecentCount)},
	{constants.TOPGAMES, make([]piratebay.Torrent, constants.RecentCount)},
	{constants.TOPVIDEOS, make([]piratebay.Torrent, constants.RecentCount)},
	{constants.TOPAPPLICATIONS, make([]piratebay.Torrent, constants.RecentCount)},
	{constants.TOPPORN, make([]piratebay.Torrent, constants.RecentCount)},
	{constants.OTHERS, make([]piratebay.Torrent, constants.RecentCount)},
}

func TestGetRecentTorrents(t *testing.T) {
	filePath := constants.LogPath
	l, err := logger.Init(filePath, logger.DebugLevel)
	assert.NoError(t, err)
	torr := NewServer(l)

	for _, rt := range recentTorrents {
		torrents, err := torr.GetRecentTorrents(rt.arg1)
		assert.NoError(t, err)
		assert.NotNil(t, torrents)
		assert.Equal(t, len(rt.expected), len(torrents))

		if len(torrents) == 1 && torrents[0].ID == "0" {
			t.Log("No torrents found")
		} else {
			str := dto.CreateResponse(torrents)
			if str == "" {
				t.Error("error creating response")
			}
		}
	}
	os.RemoveAll(path.Dir(path.Dir(filePath)))
}
