package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Write test for CreateResponse
func TestCreateResponse(t *testing.T) {
	var torrents []Torrent
	var response string

	torrents = append(torrents, Torrent{
		ID:       "1",
		InfoHash: "2",
		Name:     "3",
		NumFiles: "4",
		Size:     "5",
		Seeders:  "6",
		Leechers: "7",
	})

	response = CreateResponse(torrents)
	expected := "1) 3 - [5 B]  (6 Seeds / 7 Peers)\n\n"
	assert.Equal(t, expected, response)

}
