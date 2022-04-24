package piratebay

type Torrent struct {
	ID       string `json:"id"`
	InfoHash string `json:"info_hash"`
	Name     string `json:"name"`
	NumFiles string `json:"num_files"`
	Size     string `json:"size"`
	Seeders  string `json:"seeders"`
	Leechers string `json:"leechers"`
}
type RecentTorrent struct {
	ID       int64  `json:"id"`
	InfoHash string `json:"info_hash"`
	Name     string `json:"name"`
	NumFiles int64  `json:"num_files"`
	Size     int64  `json:"size"`
	Seeders  int64  `json:"seeders"`
	Leechers int64  `json:"leechers"`
}
