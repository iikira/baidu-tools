package baiduUtil

import "sync"

// WaitGroup 在 sync.WaitGroup 的基础上, 新增线程控制功能
type WaitGroup struct {
	wg sync.WaitGroup
	p  chan struct{}
}

// NewWaitGroup returns a pointer to a new `WaitGroup` object.
// thread 为最大并发数, 0 代表无限制
func NewWaitGroup(thread int) (w *WaitGroup) {
	w = &WaitGroup{}
	if thread <= 0 {
		return
	}
	w.p = make(chan struct{}, thread)
	return
}

// Add 在 sync.WaitGroup 的基础上, 新增线程控制功能
func (w *WaitGroup) Add(delta int) {
	w.wg.Add(delta)
	if w.p == nil {
		return
	}
	if delta >= 0 {
		for i := 0; i < delta; i++ {
			w.p <- struct{}{}
		}
	} else {
		for i := 0; i > delta; i-- {
			<-w.p
		}
	}
}

// Done 参照 sync.WaitGroup 的 Done 方法
func (w *WaitGroup) Done() {
	w.wg.Add(-1)
}

// Wait 参照 sync.WaitGroup 的 Wait 方法
func (w *WaitGroup) Wait() {
	w.wg.Wait()
}

// Thread 返回当前正在进行的任务数量
func (w *WaitGroup) Thread() int {
	return len(w.p)
}
