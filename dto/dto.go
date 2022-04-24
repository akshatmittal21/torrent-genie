package dto

import (
	"fmt"

	"github.com/akshatmittal21/torrent-genie/utils"
)

type Torrent struct {
	ID       string
	InfoHash string
	Name     string
	NumFiles string
	Size     string
	Seeders  string
	Leechers string
}

// CreateResponse : convert torrent data to text
func CreateResponse(torrents []Torrent) string {
	var response string
	for i, torrent := range torrents {
		response += fmt.Sprintf("%d) %s - [%s]  (%s Seeds / %s Peers)\n\n", i+1, torrent.Name, utils.GetFileSize(torrent.Size), torrent.Seeders, torrent.Leechers)

	}
	return response
}
