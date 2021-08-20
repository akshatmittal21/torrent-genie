package db

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/akshatmittal21/torrent-genie/constants"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	dbPath := filepath.FromSlash(constants.DBPath)

	dir := path.Dir(dbPath)

	// Creating folders if not exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatal("Unable to open database")
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
					log.Println("Error connecting to Database")
				}
			})
	}
	return dbInstance
}
