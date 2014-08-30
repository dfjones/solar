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
var _pathCache *pathCache = &pathCache{}

var counter uint64

type pathCache struct {
	paths []string
	sync.RWMutex
}

func init() {
	os.MkdirAll(dataDir, os.ModeDir|os.ModePerm)
	_pathCache.build()
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
	_pathCache.add(file.Name())
	cleanUpSignal <- struct{}{}
	return nil
}

func GetMostRecentImageFile() (*os.File, error) {
	path := _pathCache.mostRecent()
	if path == "" {
		return nil, nil
	}
	return os.Open(path)
}

func GetByIndex(i int) (*os.File, error) {
	p := _pathCache.index(i)
	if p != "" {
		return os.Open(p)
	}
	return nil, nil
}

func (p *pathCache) build() {
	p.Lock()
	defer p.Unlock()
	p.paths = listImageFiles()
}

func (p *pathCache) add(path string) {
	p.Lock()
	defer p.Unlock()
	p.paths = append(p.paths, path)
}

func (p *pathCache) index(i int) string {
	p.RLock()
	defer p.RUnlock()
	l := len(p.paths)
	if l != 0 && i < l {
		return p.paths[i]
	}
	return ""
}

func (p *pathCache) mostRecent() string {
	p.RLock()
	defer p.RUnlock()
	return p.paths[len(p.paths)-1]
}

func (p *pathCache) oldestOverLimit(limit int) []string {
	p.Lock()
	defer p.Unlock()
	diff := len(p.paths) - limit
	if diff > 0 {
		result := make([]string, diff)
		copy(result, p.paths[:diff])
		p.paths = p.paths[diff:]
		return result
	}
	return nil
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
	oldest := _pathCache.oldestOverLimit(maxFileCount)
	log.Println("Starting to remove oldest files...")
	for i, p := range oldest {
		log.Println("Removing: ", i, p)
		err := os.Remove(p)
		if err != nil {
			log.Println("Error removing file: ", err)
		}
	}
	log.Println("Finished removing oldest files")
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
