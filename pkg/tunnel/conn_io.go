package tunnel

import (
	"io"
	"sync"

	log "github.com/sirupsen/logrus"
)

// Pipe pipes data between Reader and Writer
func Pipe(src, dst io.ReadWriteCloser) (int64, int64) {
	var err error
	var sent, received int64
	var wg sync.WaitGroup
	var runOnce sync.Once
	closeFunc := func() {
		src.Close()
		dst.Close()
	}
	wg.Add(2)
	go func() {
		if received, err = io.Copy(src, dst); err != nil {
			log.Errorf("Failed to copy io from dst to src: %s", err)
		}
		runOnce.Do(closeFunc)
		wg.Done()
	}()
	go func() {
		if sent, err = io.Copy(dst, src); err != nil {
			log.Errorf("Failed to copy io from src to dst: %s", err)
		}
		runOnce.Do(closeFunc)
		wg.Done()
	}()
	wg.Wait()
	return sent, received
}
