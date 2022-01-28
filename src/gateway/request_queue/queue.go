package request_queue

import (
	"sync"
	"time"
)

type delayedRequest struct {
	delay       time.Duration
	requestFunc func() error
}

type QueueRepeater struct {
	queue   []delayedRequest
	wg      sync.WaitGroup
	quit chan bool
}

func NewQueueRepeater() *QueueRepeater {
	return &QueueRepeater{
		queue: make([]delayedRequest, 0),
		wg:    sync.WaitGroup{},
		quit:  nil,
	}
}

func (r *QueueRepeater) handleFirst() bool {
	if r.queue[0].requestFunc() == nil {
		r.queue = r.queue[1:]
		return true
	} else {
		return false
	}
}

func (r *QueueRepeater) proceed() {
	for true {
		select {
		case <-r.quit:
			r.wg.Done()
			return
		default:
			success := false
			for len(r.queue) > 0 {
				for success && len(r.queue) > 0 {
					success = r.handleFirst()
				}
				select {
				case <-time.After(r.queue[0].delay):
					success = r.handleFirst()
				}
			}
		}
	}
}

func (r *QueueRepeater) Start() {
	r.quit = make(chan bool)
	r.wg.Add(1)
	go r.proceed()
}

func (r *QueueRepeater) AddRequest(requestFunc func() error, delay time.Duration) {
	r.queue = append(r.queue, delayedRequest{
		delay:       delay,
		requestFunc: requestFunc,
	})
}

func (r *QueueRepeater) Stop() {
	r.quit <- true
	r.wg.Wait()
	close(r.quit)
}
