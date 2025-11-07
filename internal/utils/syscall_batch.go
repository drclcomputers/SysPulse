package utils

import (
	"sync"
	"time"
)

type SyscallBatcher struct {
	sync.RWMutex
	batchSize     int
	batchInterval time.Duration
	operations    []func() error
	lastBatch     time.Time
	results       []interface{}
}

func NewSyscallBatcher(batchSize int, interval time.Duration) *SyscallBatcher {
	return &SyscallBatcher{
		batchSize:     batchSize,
		batchInterval: interval,
		operations:    make([]func() error, 0),
		results:       make([]interface{}, 0),
	}
}

func (b *SyscallBatcher) Add(op func() error) {
	b.Lock()
	defer b.Unlock()

	b.operations = append(b.operations, op)

	if len(b.operations) >= b.batchSize || time.Since(b.lastBatch) >= b.batchInterval {
		b.executeBatch()
	}
}

func (b *SyscallBatcher) executeBatch() {
	if len(b.operations) == 0 {
		return
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(b.operations))

	for _, op := range b.operations {
		wg.Add(1)
		go func(operation func() error) {
			defer wg.Done()
			if err := operation(); err != nil {
				errors <- err
			}
		}(op)
	}

	wg.Wait()
	close(errors)

	b.operations = make([]func() error, 0)
	b.lastBatch = time.Now()
}

func (b *SyscallBatcher) Flush() {
	b.Lock()
	defer b.Unlock()
	b.executeBatch()
}
