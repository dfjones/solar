package image_storage

import (
	"fmt"
	"github.com/dfjones/solar/solar-server/lib/cappedlist"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync/atomic"
	"time"
)

const dataDir string = "/data/solar/images/"
const maxFileCount int = 288

var cleanUpSignal chan string = make(chan string, 10)
var pathCache *cappedlist.CappedList = cappedlist.New(maxFileCount)

var counter uint64

func init() {
	os.MkdirAll(dataDir, os.ModeDir|os.ModePerm)
	buildPathCache()
	go cleanup()
}

func buildPathCache() {
	pathCache.RegisterRemovedEntryCallback(cacheRemoveCallback)
	files := listImageFiles()
	for _, p := range files {
		pathCache.Add(p)
	}
}

func cacheRemoveCallback(entry cappedlist.Entry) {
	if e, ok := entry.(string); ok {
		cleanUpSignal <- e
	}
}

func Store(reader io.Reader) (string, error) {
	fileName := getFileName()
	file, err := os.Create(dataDir + fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return "", err
	}
	pathCache.Add(file.Name())
	return file.Name(), nil
}

func GetMostRecentImageFile() (*os.File, error) {
	path := pathCache.Last()
	if path, ok := path.(string); ok {
		if path == "" {
			return nil, nil
		}
		return os.Open(path)
	}
	return nil, nil
}

func GetByIndex(i int) (*os.File, error) {
	path := pathCache.At(i)
	if path == nil {
		return nil, nil
	}
	if path, ok := path.(string); ok {
		return os.Open(path)
	}
	return nil, nil
}

func GetAllPaths() []string {
	entries := pathCache.All()
	res := make([]string, 0)
	for _, e := range entries {
		if e, ok := e.(string); ok {
			res = append(res, e)
		}
	}
	return res
}

func getFileName() string {
	return fmt.Sprintf("%d-%d", time.Now().Unix(), getNextId())
}

func getNextId() uint64 {
	return atomic.AddUint64(&counter, uint64(1))
}

func cleanup() {
	for path := range cleanUpSignal {
		removeFile(path)
	}
}

func removeFile(path string) {
	log.Println("Removing: ", path)
	err := os.Remove(path)
	if err != nil {
		log.Println("Error removing file: ", err)
	}
	log.Println("Finished removing file", path)
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
