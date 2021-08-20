package constants

const (
	LogPath   = "./logs/system/log.log"
	DBLogPath = "./logs/db/log.log"
	DBPath    = "./db/users.db"

	TorrentCount = 7
)

const (
	ApiURL     = "https://apibay.org/q.php"
	TorrentURL = "https://itorrents.org/torrent/$$INFO_HASH$$.torrent"
)

// messages
const (
	WELCOME_MSG          = "Welcome to TorrentGenie \xF0\x9F\x98\x81	 - Get torrent magnet links \n\n Search torrents by typing a name"
	INVALID_REPLY        = "Invalid reply, please try again"
	SOMETHING_WENT_WRONG = "Something went wrong, please try again"
	NO_RESULTS           = "No torrents found"
	INVALID_COMMAND      = "Admin use only"
)

type MessageType string

const (
	Magnet  MessageType = "magnet"
	Torrent             = "torrent"
)
