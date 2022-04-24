package store

import (
	"github.com/akshatmittal21/torrent-genie/dto"
	"github.com/akshatmittal21/torrent-genie/logger"
	"github.com/akshatmittal21/torrent-genie/store/gormdb"
)

type Database interface {
	GetAllUsers() []dto.APIUser
	GetUserCount() int64
	UpsertUser(int64, dto.DBUser)
}

// New : to create single instance of the database
func New(filePath string, log logger.Logger) (Database, error) {
	return gormdb.New(filePath, log)
}
