package store

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/akshatmittal21/torrent-genie/constants"
	"github.com/akshatmittal21/torrent-genie/logger"
	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	filePath := constants.LogPath
	l, err := logger.Init(filePath, logger.DebugLevel)
	assert.NoError(t, err)

	db, err := New("./db/db.sqlite3", l)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	dbPath := filepath.FromSlash(constants.DBPath)

	dir := path.Dir(dbPath)
	os.RemoveAll(dir)
}
