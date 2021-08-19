package magnet

import (
	"net/url"
	"strings"
)

func GetLink(infoHash string, name string) string {
	magnetLink := `magnet:?xt=urn:btih:$$INFO_HASH$$&dn=$$NAME$$$$TRACKERS$$`
	name = url.QueryEscape(name)
	magnetLink = strings.Replace(magnetLink, "$$INFO_HASH$$", infoHash, 1)
	magnetLink = strings.Replace(magnetLink, "$$NAME$$", name, 1)
	magnetLink = strings.Replace(magnetLink, "$$TRACKERS$$", printTrackers(), 1)
	return magnetLink
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
