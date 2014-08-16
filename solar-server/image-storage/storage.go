package image_storage

import (
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var dataDir string = "/data/solar/images/"
var counter uint64
var mostRecentImageInfo MostRecentImageInfo

type MostRecentImageInfo struct {
	path string
	sync.RWMutex
}

func init() {
	os.MkdirAll(dataDir, os.ModeDir|os.ModePerm)
}

func Store(reader io.Reader) error {
	fileName := getFileName()
	file, err := os.Create(dataDir + fileName)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}
	updateMostRecent(file)
	return nil
}

func GetMostRecentImageReader() (io.Reader, error) {
	path := getMostRecentPath()
	if path == "" {
		return nil, nil
	}
	return os.Open(path)
}

func getMostRecentPath() string {
	mostRecentImageInfo.RLock()
	defer mostRecentImageInfo.RUnlock()
	return mostRecentImageInfo.path
}

func updateMostRecent(file *os.File) {
	mostRecentImageInfo.Lock()
	defer mostRecentImageInfo.Unlock()
	mostRecentImageInfo.path = file.Name()
}

func getFileName() string {
	return fmt.Sprintf("%d-%d", time.Now().Unix(), getNextId())
}

func getNextId() uint64 {
	return atomic.AddUint64(&counter, uint64(1))
}
