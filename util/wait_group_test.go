package baiduUtil

import (
	"fmt"
	"testing"
	"time"
)

func TestWg(t *testing.T) {
	wg := NewWaitGroup(2)
	for i := 0; i < 60; i++ {
		wg.Add()
		go func() {
			fmt.Println(i, wg.Thread())
			time.Sleep(1e9)
			wg.Done()
		}()
	}
	wg.Wait()
}
