package processes

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/process"
)

type ProcessCache struct {
	sync.RWMutex
	processes      map[int32]*ProcessInfo
	lastUpdate     time.Time
	updateInterval time.Duration
}

type ProcessInfo struct {
	Name       string
	CPU        float64
	Memory     float32
	Status     string
	Username   string
	CmdLine    string
	CreateTime int64
	NumThreads int32
	MemoryInfo *process.MemoryInfoStat
	LastUpdate time.Time
	TTL        time.Duration
}

var (
	defaultCache *ProcessCache
	once         sync.Once
)

func GetProcessCache() *ProcessCache {
	once.Do(func() {
		defaultCache = &ProcessCache{
			processes:      make(map[int32]*ProcessInfo),
			updateInterval: 2 * time.Second,
		}
	})
	return defaultCache
}

func (pc *ProcessCache) Get(pid int32) (*ProcessInfo, bool) {
	pc.RLock()
	defer pc.RUnlock()

	info, exists := pc.processes[pid]
	if !exists {
		return nil, false
	}

	if time.Since(info.LastUpdate) > info.TTL {
		return nil, false
	}

	return info, true
}

func (pc *ProcessCache) Set(pid int32, info *ProcessInfo) {
	pc.Lock()
	defer pc.Unlock()

	pc.processes[pid] = info
}

func (pc *ProcessCache) Clear() {
	pc.Lock()
	defer pc.Unlock()

	pc.processes = make(map[int32]*ProcessInfo)
}
