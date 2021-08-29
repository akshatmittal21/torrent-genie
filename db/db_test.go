package db

import (
	"os"
	"path"
	"path/filepath"
	"sync"
	"testing"

	"github.com/akshatmittal21/torrent-genie/constants"
	"gorm.io/gorm"
)

func TestGetInstance(t *testing.T) {
	instance := GetInstance()
	if instance == nil {
		t.Error("GetInstance() failed expected pointer got nil")
	}
	dbPath := filepath.FromSlash(constants.DBPath)

	dir := path.Dir(dbPath)
	os.RemoveAll(dir)
}

func TestSingleInstance(t *testing.T) {

	instanceCh := make(chan *gorm.DB, 10)
	instanceArr := make([]*gorm.DB, 0)
	var wg sync.WaitGroup
	go func() {
		for i := range instanceCh {
			instanceArr = append(instanceArr, i)
		}
	}()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			instance := GetInstance()
			instanceCh <- instance
		}()
	}
	wg.Wait()
	close(instanceCh)
	tempInstance := instanceArr[0]
	for i := 1; i < len(instanceArr); i++ {
		if instanceArr[i] == nil {
			t.Error("GetInstance() failed expected pointer got nil")
		}
		if tempInstance != instanceArr[i] {
			t.Errorf("GetInstance() failed expected same pointer (%v) got different (%v)", tempInstance, instanceArr[i])
		}
		tempInstance = instanceArr[i]
	}
	dbPath := filepath.FromSlash(constants.DBPath)

	dir := path.Dir(dbPath)
	os.RemoveAll(dir)

}
