package db

import (
	oslog "log"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/akshatmittal21/torrent-genie/constants"
	"github.com/akshatmittal21/torrent-genie/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

type UserConfig struct {
	UserID    int64 `gorm:"index,primary_key"`
	FirstName string
	LastName  string
	UserName  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

// GetInstance : to create single instance of the database
func GetInstance() *gorm.DB {
	newLogger := gormlogger.New(
		oslog.New(os.Stdout, "\r\n", oslog.LstdFlags), // io writer
		gormlogger.Config{
			SlowThreshold:             time.Second,       // Slow SQL threshold
			LogLevel:                  gormlogger.Silent, // Log level
			IgnoreRecordNotFoundError: true,              // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,             // Disable color
		},
	)

	dbPath := filepath.FromSlash(constants.DBPath)

	dir := path.Dir(dbPath)

	// Creating folders if not exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			oslog.Fatal("Unable to open database")
		}
	}
	if dbInstance == nil {
		once.Do(
			func() {
				var err error
				// log.Println("Creating Single Instance Now")
				dbInstance, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{Logger: newLogger})

				// Creating DB schemas
				dbInstance.AutoMigrate(&UserConfig{})

				if err != nil {
					logger.Error("Error connecting to Database")
				}
			})
	}
	return dbInstance
}
