package usecase

import (
	"errors"
	"time"
)

const (
	open     = 0
	halfOpen = 1
	closed   = 2
)

type CircuitBreaker struct {
	FailsMax int
	fails    int
	active   bool
	state    int
}

func NewCircuitBreaker(failsMax int) *CircuitBreaker {
	return &CircuitBreaker{
		FailsMax: failsMax,
		fails:    0,
		active:   false,
		state:    closed,
	}
}

func (cb *CircuitBreaker) Call(requestFunc func() (interface{}, error)) (interface{}, error) {
	if cb.state != open {
		res, err := requestFunc()
		if err == nil {
			cb.fails = 0
			if cb.state == halfOpen {
				cb.state = closed
			}
		} else {
			if cb.state == halfOpen {
				cb.state = open
			} else {
				cb.fails += 1
				if cb.fails >= cb.FailsMax {
					cb.state = open
				}
			}
		}
		return res, err
	} else {
		if !cb.active {
			cb.active = true
			select {
			case <-time.After(2 * time.Second):
				cb.state = halfOpen
				cb.active = false
			}
		}
		return nil, errors.New("Circuit breaker is open")
	}
}
