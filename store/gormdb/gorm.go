package gormdb

import (
	oslog "log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/akshatmittal21/torrent-genie/dto"
	"github.com/akshatmittal21/torrent-genie/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type DB struct {
	db *gorm.DB
}

func New(filePath string, log logger.Logger) (*DB, error) {
	newLogger := gormlogger.New(
		oslog.New(os.Stdout, "\r\n", oslog.LstdFlags), // io writer
		gormlogger.Config{
			SlowThreshold:             time.Second,       // Slow SQL threshold
			LogLevel:                  gormlogger.Silent, // Log level
			IgnoreRecordNotFoundError: true,              // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,             // Disable color
		},
	)

	dbPath := filepath.FromSlash(filePath)

	dir := path.Dir(dbPath)

	// Creating folders if not exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			oslog.Fatal("Unable to open database")
		}
	}

	var err error
	// log.Println("Creating Single Instance Now")
	dbInstance, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Error("Error connecting to Database %w", err)
	}
	// Creating DB schemas
	err = dbInstance.AutoMigrate(&userConfig{})
	if err != nil {
		log.Error("Error migrating database %w", err)
	}
	return &DB{db: dbInstance}, nil

}

// get all users
func (db *DB) GetAllUsers() []dto.APIUser {
	var users []dto.APIUser
	db.db.Model(&userConfig{}).Find(&users)
	return users
}

func (db *DB) UpsertUser(userID int64, user dto.DBUser) {
	dbUser := userConfig{
		UserID:    userID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UserName:  user.UserName,
	}
	rowsAffected := db.db.Model(&userConfig{}).Where("user_id = ?", userID).Updates(&dbUser).RowsAffected
	if rowsAffected == 0 {
		db.db.Create(&dbUser)
	}
}

func (db *DB) GetUserCount() int64 {
	var count int64
	db.db.Find(&userConfig{}).Count(&count)
	return count
}
