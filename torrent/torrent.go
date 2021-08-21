package torrent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/akshatmittal21/torrent-genie/constants"
	"github.com/akshatmittal21/torrent-genie/logger"
)

type Torrent struct {
	ID       string `json:"id"`
	InfoHash string `json:"info_hash"`
	Name     string `json:"name"`
	NumFiles string `json:"num_files"`
	Size     string `json:"size"`
	Seeders  string `json:"seeders"`
	Leechers string `json:"leechers"`
}

// GetTorrents : to get torrents from apibay
func GetTorrents(searchText string) []Torrent {
	var torrents []Torrent

	u, err := url.Parse(constants.ApiURL)
	if err != nil {
		fmt.Println(err)
		return torrents
	}
	q := u.Query()
	searchText = url.QueryEscape(searchText)
	q.Set("q", searchText)
	q.Set("cat", "")
	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		logger.Error("GetTorrents: fetch err ", err)
		return torrents
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("GetTorrents: reading err ", err)
		return torrents
	}
	err = json.Unmarshal(body, &torrents)
	if err != nil {
		logger.Error("GetTorrents: unmarshal err ", err)
		return torrents
	}
	if len(torrents) <= constants.TorrentCount {
		return torrents
	}
	// fmt.Println(torrents)
	return torrents[:constants.TorrentCount]
}

// CreateResponse : convert torrent data to text
func CreateResponse(torrents []Torrent) string {
	var response string
	for i, torrent := range torrents {
		response += fmt.Sprintf("%d) %s - [%s]  (%sSeeds / %s Peers)\n\n", i+1, torrent.Name, getSize(torrent.Size), torrent.Seeders, torrent.Leechers)

	}
	return response
}
