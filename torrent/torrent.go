package torrent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Torrent struct {
	ID       string `json:"id"`
	InfoHash string `json:"info_hash"`
	Name     string `json:"name"`
	NumFiles string `json:"num_files"`
	Size     string `json:"size"`
}

const torrentCount = 5
const apiURL = "https://apibay.org/q.php"

// GetTorrents : to get torrents from apibay
func GetTorrents(searchText string) []Torrent {
	var torrents []Torrent

	u, err := url.Parse(apiURL)
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
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(body, &torrents)
	if err != nil {
		log.Fatal(err)
	}
	if len(torrents) <= torrentCount {
		return torrents
	}
	// fmt.Println(torrents)
	return torrents[:torrentCount]
}

// CreateResponse : convert torrent data to text
func CreateResponse(torrents []Torrent) string {
	var response string
	for i, torrent := range torrents {
		response += fmt.Sprintf("%d : %s - (%s)\n\n", i+1, torrent.Name, getSize(torrent.Size))

	}
	return response
}
