package torrent

import (
	"os"
	"path"
	"testing"

	"github.com/akshatmittal21/torrent-genie/constants"
)

func TestGetTorrent(t *testing.T) {
	torrents := GetTorrents("abcd")
	if len(torrents) == 0 {
		t.Error("error fetching torrents")
	}

	torrents = GetTorrents("oquu2gdbi")
	if len(torrents) == 0 {
		t.Error("error fetching torrents")
	}
	if torrents[0].ID == "0" {
		t.Log("No torrents found")
	}

	torrents = GetTorrents("")
	if len(torrents) == 0 {
		t.Error("error fetching torrents")
	}
	if torrents[0].ID == "0" {
		t.Log("No torrents found")
	}
}

type recentTorrent struct {
	arg1     constants.RecentCode
	expected []Torrent
}

var recentTorrents = []recentTorrent{
	{constants.RecentAllCode, make([]Torrent, 10)},
	{constants.AudioCode, make([]Torrent, 10)},
	{constants.GamesCode, make([]Torrent, 10)},
	{constants.VideoCode, make([]Torrent, 10)},
	{constants.ApplicationsCode, make([]Torrent, 10)},
	{constants.PornCode, make([]Torrent, 10)},
	{constants.OthersCode, make([]Torrent, 10)},
}

func TestGetRecentTorrents(t *testing.T) {
	for _, rt := range recentTorrents {
		torrents := GetRecentTorrents(rt.arg1)
		if torrents == nil {
			t.Error("error fetching torrents")
		}
		if len(torrents) != len(rt.expected) {
			t.Errorf("torrents length as not expected for (%v) code", rt.arg1)
		}
		if len(torrents) == 1 && torrents[0].ID == "0" {
			t.Log("No torrents found")
		} else {
			str := CreateResponse(torrents)
			if str == "" {
				t.Error("error creating response")
			}
		}
	}
	filePath := constants.LogPath
	os.RemoveAll(path.Dir(path.Dir(filePath)))
}
