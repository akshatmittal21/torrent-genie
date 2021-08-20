package magnet

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/akshatmittal21/torrent-genie/constants"
	"github.com/akshatmittal21/torrent-genie/logger"
)

func GetLink(infoHash string, name string) string {
	magnetLink := `magnet:?xt=urn:btih:$$INFO_HASH$$&dn=$$NAME$$$$TRACKERS$$`
	name = url.QueryEscape(name)
	magnetLink = strings.Replace(magnetLink, "$$INFO_HASH$$", infoHash, 1)
	magnetLink = strings.Replace(magnetLink, "$$NAME$$", name, 1)
	magnetLink = strings.Replace(magnetLink, "$$TRACKERS$$", printTrackers(), 1)
	return magnetLink
}

func GetFile(infoHash string) []byte {
	torrentLink := constants.TorrentURL
	torrentLink = strings.Replace(torrentLink, "$$INFO_HASH$$", infoHash, 1)
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(torrentLink)
	if err != nil {
		logger.Error("GetFile: fetch err ", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("GetFile: reading err ", err)
		return nil
	}
	return body
}

func printTrackers() string {
	tr := "&tr=" + url.QueryEscape("udp://tracker.coppersurfer.tk:6969/announce")
	tr += "&tr=" + url.QueryEscape("udp://tracker.openbittorrent.com:6969/announce")
	tr += "&tr=" + url.QueryEscape("udp://9.rarbg.to:2710/announce")
	tr += "&tr=" + url.QueryEscape("udp://9.rarbg.me:2780/announce")
	tr += "&tr=" + url.QueryEscape("udp://9.rarbg.to:2730/announce")
	tr += "&tr=" + url.QueryEscape("udp://tracker.opentrackr.org:1337")
	tr += "&tr=" + url.QueryEscape("udp://movies.zsw.ca:6969/announce")
	tr += "&tr=" + url.QueryEscape("udp://tracker.dler.org:6969/announce")
	tr += "&tr=" + url.QueryEscape("udp://open.stealth.si:80/announce")
	tr += "&tr=" + url.QueryEscape("udp://tracker.0x.tf:6969/announce")
	return tr
}
