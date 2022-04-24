package torrent

import (
	"github.com/akshatmittal21/torrent-genie/constants"
	"github.com/akshatmittal21/torrent-genie/dto"
	"github.com/akshatmittal21/torrent-genie/logger"
	"github.com/akshatmittal21/torrent-genie/torrent/piratebay"
)

type Server interface {
	GetRecentTorrents(constants.Command) ([]dto.Torrent, error)
	GetTorrents(string) ([]dto.Torrent, error)
	IsServerAlive() bool
}

func NewServer(log logger.Logger) Server {
	return piratebay.NewServer(log)
}
