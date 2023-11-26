package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrNoTasksSupplied = errors.New("no tasks supplied for processing")
var ErrZeroWorkersCount = errors.New("supplied workers counter is zero")

type Task func() error

func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return ErrNoTasksSupplied
	}
	if n == 0 {
		return ErrZeroWorkersCount
	}
	queue := make(chan Task, len(tasks))
	done := make(chan struct{})

	go toQueue(done, tasks, queue)
	out := merge(queue, n, done)
	// Wait until all remaining goroutines finish to mitigate race conditions
	defer func() {
		for range out {

		}
	}()

	var skiperrors bool
	if m == 0 {
		skiperrors = true
	}
	for o := range out {
		select {
		case <-done:
			return nil
		default:
			if skiperrors == true {
				continue
			}
			if o != nil {
				m--
			}
			if m <= 0 {
				close(done)
				return ErrErrorsLimitExceeded
			}
		}
	}

	return nil
}

func toQueue(done <-chan struct{}, tasks []Task, queue chan<- Task) {
	defer close(queue)
	for _, t := range tasks {
		select {
		case <-done:
			return
		default:
			queue <- t
		}
	}
}

func worker(queue <-chan Task, done <-chan struct{}, out chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range queue {
		select {
		case <-done:
			return
		default:
			out <- t()
		}
	}
}

func merge(queue <-chan Task, n int, done chan struct{}) <-chan error {
	wg := &sync.WaitGroup{}
	out := make(chan error)
	for i := 0; i < n; i++ {
		select {
		case <-done:
			break
		default:
			wg.Add(1)
			go worker(queue, done, out, wg)
		}
	}
	go func() {
		defer close(out)
		wg.Wait()

		select {
		case <-done:
			return
		default:
			close(done)
			return
		}
	}()

	return out
}
