package constants

// Internal config
const (
	LogPath   = "./logs/system/log.log"
	DBLogPath = "./logs/db/log.log"
	DBPath    = "./db/users.db"

	TorrentCount = 15
	RecentCount  = 15
)

// URLs

const (
	ApiURL           = "https://apibay.org/q.php"
	RecentTorrentURL = "https://apibay.org/precompiled/data_top100_${CODE}.json"
	TorrentURL       = "https://itorrents.org/torrent/${INFO_HASH}.torrent"
)

const (
	MagnetLink string = "magnet:?xt=urn:btih:${INFO_HASH}&dn=${NAME}${TRACKERS}"
)

// Recommendation
var (
	IsRecommendationOn bool = true
)

const (
	RecommendMsg string = `Hi ${NAME} ✌️

	Thank you for letting me save your time ❤️.
	If you like me, why don't invite your friends too. I am sure they will like me too 😁.
	Please share the below link to let them know about me 😊
	https://t.me/torrent_genie_bot
	
	Ciao 🤖`
)

// Recent codes
type RecentCode string

const (
	RecentAllCode    RecentCode = "recent"
	AudioCode                   = "100"
	GamesCode                   = "400"
	VideoCode                   = "200"
	PornCode                    = "500"
	ApplicationsCode            = "300"
	OthersCode                  = "600"
)

//commands
const (
	START           string = "start"
	USERS                  = "users"
	RECENT                 = "recent"
	TOPVIDEOS              = "topvideos"
	TOPAUDIO               = "topaudio"
	TOPGAMES               = "topgames"
	TOPPORN                = "topporn"
	TOPAPPLICATIONS        = "topapplications"
	OTHERS                 = "others"
)

// messages
const (
	WELCOME_MSG          = "Welcome to TorrentGenie \xF0\x9F\x98\x81	 - Get torrent files | magnet links \n\n Start by typing a name"
	INVALID_REPLY        = "Invalid reply, please try again"
	SOMETHING_WENT_WRONG = "Something went wrong, please try again"
	NO_RESULTS           = "No torrents found"
	INVALID_COMMAND      = "Invalid Command"
	ADMIN_COMMAND        = "Admin use only"
	SERVER_DOWN          = "Torrent server is down, please try again later"
)

type MessageType string

const (
	Magnet  MessageType = "magnet"
	Torrent             = "torrent"
)

// ping timeout
const (
	PingTimeout = 10 * 60
)
