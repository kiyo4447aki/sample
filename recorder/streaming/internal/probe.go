// internal/probe.go
package internal

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-gst/go-gst/gst"
)

func StartProbe(
	ctx context.Context,
	appsink *gst.Element,
	timeout time.Duration,
	interval time.Duration,
	onStall func(),
) error {
	var mu sync.Mutex
	last := time.Now()

	_, err := appsink.Connect("new-sample", func(_ *gst.Element) gst.FlowReturn {
		mu.Lock()
		last = time.Now()
		mu.Unlock()
		return gst.FlowOK
	})
	if err != nil {
		return fmt.Errorf("new-sampleシグナルのコールバックを設定できませんでした：%w", err)}

	// 監視用 goroutine
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				mu.Lock()
				elapsed := time.Since(last)
				mu.Unlock()
				if elapsed > timeout {
					onStall()
				}
			}
		}
	}()

	return nil
}
