package server_info

import (
	"cms/internal/controller/http/response"
	"cms/pkg/logger"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"runtime"
	"sort"
	"time"

	"github.com/gorilla/mux"
)

const RouteName = "api.admin.server_info"

const bytesInKilobyte = 1024

type (
	responseData struct {
		Host    string      `json:"Host"`
		Runtime runtimeInfo `json:"Runtime"`
		Files   []fileInfo  `json:"Files"`
	}
	runtimeInfo struct {
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
)

type handler struct {
	files     fs.FS
	startTime time.Time
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	host, err := os.Hostname()
	if err != nil {
		logger.Load(r.Context()).Error(err, "failed to get hostname")
		response.InternalServerError(rw, r)
		return
	}

	memory := &runtime.MemStats{} //nolint:exhaustruct
	runtime.ReadMemStats(memory)

	info := runtimeInfo{
		NumCPU:       runtime.NumCPU(),
		Memory:       memory.Alloc,
		MemSys:       memory.Sys / bytesInKilobyte,
		HeapAlloc:    memory.HeapAlloc / bytesInKilobyte,
		HeapSys:      memory.HeapSys / bytesInKilobyte,
		HeapIdle:     memory.HeapIdle / bytesInKilobyte,
		HeapInuse:    memory.HeapInuse / bytesInKilobyte,
		HeapReleased: memory.HeapReleased / bytesInKilobyte,
		NextGC:       memory.NextGC / bytesInKilobyte,
		Goroutines:   runtime.NumGoroutine(),
		UpTime:       uint64(time.Since(h.startTime).Seconds()),
		Time:         time.Now().Format(time.RFC1123Z),
	}

	var files fileInfosSortedByName
	if err := fs.WalkDir(h.files, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return fmt.Errorf("get info: %w", err)
		}
		files = append(files, fileInfo{
			Name: path,
			Size: info.Size(),
		})
		return nil
	}); err != nil {
		logger.Load(r.Context()).Error(err, "failed to walk to embed fs")
		response.InternalServerError(rw, r)
		return
	}
	sort.Sort(files)

	response.JSON(rw, r, http.StatusOK, &responseData{
		Host:    host,
		Runtime: info,
		Files:   files,
	})
}

func newHandler(files fs.FS, startTime time.Time) *handler {
	return &handler{
		files:     files,
		startTime: startTime,
	}
}

func Setup(router *mux.Router, files fs.FS, startTime time.Time) {
	router.
		Name(RouteName).
		Methods(http.MethodGet).
		Path("/server-info").
		Handler(newHandler(files, startTime))
}
