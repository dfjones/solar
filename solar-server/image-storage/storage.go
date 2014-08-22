package image_storage

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

const dataDir string = "/data/solar/images/"
const maxFileCount int = 288

var cleanUpSignal chan struct{} = make(chan struct{})
var mostRecentImageInfo MostRecentImageInfo

var counter uint64

type MostRecentImageInfo struct {
	path string
	sync.RWMutex
}

func init() {
	os.MkdirAll(dataDir, os.ModeDir|os.ModePerm)
	imageFiles := listImageFiles()
	if len(imageFiles) > 0 {
		mostRecentImageInfo.Lock()
		defer mostRecentImageInfo.Unlock()
		mostRecentImageInfo.path = imageFiles[len(imageFiles)-1]
	}
	go cleanup()
}

func Store(reader io.Reader) error {
	fileName := getFileName()
	file, err := os.Create(dataDir + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}
	updateMostRecent(file)
	cleanUpSignal <- struct{}{}
	return nil
}

func GetMostRecentImageFile() (*os.File, error) {
	path := getMostRecentPath()
	if path == "" {
		return nil, nil
	}
	return os.Open(path)
}

func GetByIndex(i int) (*os.File, error) {
	images := listImageFiles()
	l := len(images)
	if l == 0 || i >= l {
		return nil, nil
	}
	return os.Open(images[i])
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

func cleanup() {
	for _ = range cleanUpSignal {
		// walk all the files in our data directory and gather them into a slice
		removeOldestFiles()
	}
}

func removeOldestFiles() {
	imageFiles := listImageFiles()
	extraFileCount := len(imageFiles) - maxFileCount
	log.Println("Starting to remove oldest files...")
	for i, p := range imageFiles {
		if i >= extraFileCount {
			return
		}
		log.Println("Removing: ", i, p)
		err := os.Remove(p)
		if err != nil {
			log.Println("Error removing file: ", err)
		}
	}
}

func listImageFiles() []string {
	imageFiles := make([]string, 0)
	filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if path != dataDir {
			imageFiles = append(imageFiles, path)
		}
		return nil
	})

	// the files use unix timestamps as their names, so sorting them like this should order
	// them oldest to newest
	sort.Strings(imageFiles)
	return imageFiles
}
