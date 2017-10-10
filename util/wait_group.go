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

// AddDelta 在 sync.WaitGroup 的基础上, 新增线程控制功能
func (w *WaitGroup) AddDelta() {
	w.wg.Add(1)
	if w.p == nil {
		return
	}
	w.p <- struct{}{}
}

// Done 在 sync.WaitGroup 的基础上, 新增线程控制功能
func (w *WaitGroup) Done() {
	w.wg.Done()
	if w.p == nil {
		return
	}
	<-w.p
}

// Wait 参照 sync.WaitGroup 的 Wait 方法
func (w *WaitGroup) Wait() {
	w.wg.Wait()
}

// Thread 返回当前正在进行的任务数量
func (w *WaitGroup) Thread() int {
	return len(w.p)
}
