package handlers

import (
	"fmt"
	"github.com/go-logr/logr"
	"io/fs"
	"net/http"
	"os"
	"{{shortName}}/assets"
	"{{shortName}}/core/handler"
	"{{shortName}}/core/sessions"
	"runtime"
	"sort"
	"time"
)

var (
	startTime = time.Now()
)

type info struct {
	Host    string       `json:"Host"`
	Runtime *runtimeInfo `json:"Runtime"`
	Files   []fileInfo   `json:"Files"`
}

// runtimeInfo defines runtime part of service information
type runtimeInfo struct {
	NumCPU       int    `json:"NumCPU"`
	Memory       uint64 `json:"Memory"`
	MemSys       uint64 `json:"MemSys"`
	HeapAlloc    uint64 `json:"HeapAlloc"`
	HeapSys      uint64 `json:"HeapSys"`
	HeapIdle     uint64 `json:"HeapIdle"`
	HeapInuse    uint64 `json:"HeapInuse"`
	HeapReleased uint64 `json:"HeapRealease"`
	NextGC       uint64 `json:"NextGC"`
	Goroutines   int    `json:"Goroutines"`
	UpTime       uint64 `json:"UpTime"`
	Time         string `json:"Time"`
}

type fileInfo struct {
	Name string `json:"Name"`
	Size int64  `string:"Size"`
}

type fileInfosSortedByName []fileInfo

func (fi fileInfosSortedByName) Len() int           { return len(fi) }
func (fi fileInfosSortedByName) Less(i, j int) bool { return fi[i].Name < fi[j].Name }
func (fi fileInfosSortedByName) Swap(i, j int)      { fi[i], fi[j] = fi[j], fi[i] }

func infoHandler(_ *handler.Context, _ interface{}) (interface{}, error) {
	host, _ := os.Hostname()
	memory := &runtime.MemStats{}
	runtime.ReadMemStats(memory)
	rt := &runtimeInfo{
		NumCPU:       runtime.NumCPU(),
		Memory:       memory.Alloc,
		MemSys:       memory.Sys / 1024,
		HeapAlloc:    memory.HeapAlloc / 1024,
		HeapSys:      memory.HeapSys / 1024,
		HeapIdle:     memory.HeapIdle / 1024,
		HeapInuse:    memory.HeapInuse / 1024,
		HeapReleased: memory.HeapReleased / 1024,
		NextGC:       memory.NextGC / 1024,
		Goroutines:   runtime.NumGoroutine(),
		UpTime:       uint64(time.Since(startTime).Seconds()),
		Time:         time.Now().Format(time.RFC1123Z),
	}

	var files fileInfosSortedByName
	fsys, err := fs.Sub(assets.Assets, "files")
	if nil != err {
		return nil, fmt.Errorf("can't get fs.sub: %w", err)
	}
	err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if nil != err {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if nil != err {
			return fmt.Errorf("get info: %w", err)
		}
		files = append(files, fileInfo{
			Name: path,
			Size: info.Size(),
		})
		return nil
	})
	if nil != err {
		return nil, fmt.Errorf("can't walk to embed fs: %w", err)
	}

	sort.Sort(files)

	return &info{
		Host:    host,
		Runtime: rt,
		Files:   files,
	}, nil
}

func NewInfoHandler(logger logr.Logger, sessionStore *sessions.Store) http.HandlerFunc {
	return handler.CreateFromProcessFunc(logger, infoHandler, handler.WithOnlyForSuperuser(sessionStore))
}
