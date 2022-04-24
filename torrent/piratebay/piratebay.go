package piratebay

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/akshatmittal21/torrent-genie/constants"
	"github.com/akshatmittal21/torrent-genie/dto"
	"github.com/akshatmittal21/torrent-genie/logger"
)

type PirateBay struct {
	log logger.Logger
}

func NewServer(log logger.Logger) PirateBay {
	return PirateBay{log: log}
}

func (p PirateBay) GetRecentTorrents(command constants.Command) ([]dto.Torrent, error) {
	var (
		torrents []dto.Torrent
		err      error
	)
	switch command {
	case constants.RECENT:
		torrents, err = p.getRecentTorrents(RecentAllCode)
	case constants.TOPAUDIO:
		torrents, err = p.getRecentTorrents(AudioCode)
	case constants.TOPVIDEOS:
		torrents, err = p.getRecentTorrents(VideoCode)
	case constants.TOPAPPLICATIONS:
		torrents, err = p.getRecentTorrents(ApplicationsCode)
	case constants.TOPPORN:
		torrents, err = p.getRecentTorrents(PornCode)
	case constants.TOPGAMES:
		torrents, err = p.getRecentTorrents(GamesCode)
	case constants.OTHERS:
		torrents, err = p.getRecentTorrents(OthersCode)
	}
	return torrents, err

}
func (p PirateBay) getRecentTorrents(code RecentCode) ([]dto.Torrent, error) {
	var recentTorrents []RecentTorrent
	recentURL := strings.Replace(constants.RecentTorrentURL, "${CODE}", string(code), 1)
	resp, err := http.Get(recentURL)
	if err != nil {
		err = fmt.Errorf("GetTorrents: fetch err %w", err)
		p.log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("GetTorrents: read err %w", err)
		p.log.Error(err)
		return nil, err
	}
	err = json.Unmarshal(body, &recentTorrents)
	if err != nil {
		err = fmt.Errorf("GetTorrents: unmarshal err %w", err)
		p.log.Error(err)
		return nil, err
	}
	torrents := make([]dto.Torrent, len(recentTorrents))
	for i, torrent := range recentTorrents {
		torrents[i].ID = strconv.FormatInt(torrent.ID, 10)
		torrents[i].InfoHash = torrent.InfoHash
		torrents[i].Name = torrent.Name
		torrents[i].NumFiles = strconv.FormatInt(torrent.NumFiles, 10)
		torrents[i].Size = strconv.FormatInt(torrent.Size, 10)
		torrents[i].Seeders = strconv.FormatInt(torrent.Seeders, 10)
		torrents[i].Leechers = strconv.FormatInt(torrent.Leechers, 10)
	}
	if len(torrents) <= constants.RecentCount {
		return torrents, nil
	}
	// fmt.Println(torrents)
	return torrents[:constants.RecentCount], nil
}

// GetTorrents : to get torrents from apibay
func (p PirateBay) GetTorrents(searchText string) ([]dto.Torrent, error) {
	var torrents []Torrent
	u, err := url.Parse(constants.ApiURL)
	if err != nil {
		err = fmt.Errorf("GetTorrents: parse err %w", err)
		p.log.Error(err)
		return nil, err
	}
	q := u.Query()
	searchText = url.QueryEscape(searchText)
	q.Set("q", searchText)
	q.Set("cat", "")
	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		err = fmt.Errorf("GetTorrents: fetch err %w", err)
		p.log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("GetTorrents: read err %w", err)
		p.log.Error(err)
		return nil, err
	}
	err = json.Unmarshal(body, &torrents)
	if err != nil {
		err = fmt.Errorf("GetTorrents: unmarshal err %w", err)
		p.log.Error(err)
		return nil, err
	}
	torrentsDTO := make([]dto.Torrent, len(torrents))
	for i, torrent := range torrents {
		torrentsDTO[i].ID = torrent.ID
		torrentsDTO[i].InfoHash = torrent.InfoHash
		torrentsDTO[i].Name = torrent.Name
		torrentsDTO[i].NumFiles = torrent.NumFiles
		torrentsDTO[i].Size = torrent.Size
		torrentsDTO[i].Seeders = torrent.Seeders
		torrentsDTO[i].Leechers = torrent.Leechers
	}
	if len(torrents) <= constants.TorrentCount {
		return torrentsDTO, nil
	}
	// fmt.Println(torrents)
	return torrentsDTO[:constants.TorrentCount], nil
}

func (p PirateBay) IsServerAlive() bool {
	resp, err := http.Get(constants.ApiURL)
	if err != nil {
		p.log.Error("ServerStatus: fetch err ", err)
		return false
	}
	if resp != nil && resp.StatusCode == 200 {
		return true
	}
	p.log.Error("ServerStatus: server is down ")
	return false

}
