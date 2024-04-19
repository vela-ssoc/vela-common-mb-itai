package gopool

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Func func()

func (f Func) Run() { f() }

type Runner interface {
	Run()
}

type Executor interface {
	Submit(Runner)
	Execute(func())
	Close()
}

// New 创建协程池。
func New(size, buff int, idle time.Duration) Executor {
	// 确保协程数必须大于 0
	if size <= 0 {
		size = runtime.NumCPU() * 4
		if size <= 0 {
			size = 8
		}
	}
	if buff < 0 {
		buff = 0
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &bufferedPool{
		size:   size,
		idle:   idle,
		queue:  make(chan Runner, buff),
		ctx:    ctx,
		cancel: cancel,
	}
}

type bufferedPool struct {
	mutex sync.Mutex

	// size 协程数最大个数
	size int

	// work 正在运行的协程数
	work int

	// idle 空闲时间，当协程执行完 Runner 后会从任务队列里面拿 Runner。
	// 如果超过空闲时间仍未拿到新的 Runner 则协程退出运行。
	idle time.Duration

	// queue 存放 Runner 队列
	queue chan Runner

	// closed 是否已经关闭
	closed atomic.Bool
	ctx    context.Context
	cancel context.CancelFunc
}

func (bp *bufferedPool) Submit(r Runner) {
	if r == nil || bp.closed.Load() {
		return
	}

	// 先判断是否可以新建协程执行，
	// 可以创建就创建协程；
	// 协程已满则将任务放到任务队列，等待空闲协程去拿取并执行。
	if bp.acquire() {
		go bp.worker(r)
		return
	}

	select {
	case <-bp.ctx.Done():
	case bp.queue <- r:
	}
}

func (bp *bufferedPool) Execute(f func()) {
	if f != nil {
		bp.Submit(Func(f))
	}
}

func (bp *bufferedPool) Close() {
	if bp.closed.CompareAndSwap(false, true) {
		bp.cancel()
	}
}

func (bp *bufferedPool) worker(r Runner) {
	defer bp.release()

	bp.execute(r)

	idle := bp.idle
	if idle <= 0 {
		return
	}

	timer := time.NewTimer(bp.idle)
	defer timer.Stop()

	var over bool
	for !over {
		select {
		case <-timer.C:
			over = true
		case rn := <-bp.queue:
			bp.execute(rn)
			timer.Reset(idle)
		}
	}
}

func (bp *bufferedPool) acquire() bool {
	bp.mutex.Lock()
	ok := bp.size > bp.work
	if ok {
		bp.work++
	}
	bp.mutex.Unlock()

	return ok
}

func (bp *bufferedPool) release() {
	bp.mutex.Lock()
	bp.work--
	bp.mutex.Unlock()
}

func (bp *bufferedPool) execute(r Runner) {
	defer func() { recover() }()
	r.Run()
}
