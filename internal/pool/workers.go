package pool

import (
	"context"
	"log"
	"sync"
)

type Task func()

type Pool struct {
	tasks  chan Task
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewPool(size int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Pool{
		tasks:  make(chan Task, 100),
		ctx:    ctx,
		cancel: cancel,
	}

	for i := 0; i < size; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}

	return p
}

func (p *Pool) worker(id int) {
	defer p.wg.Done()

	log.Printf("[worker %d] started", id)
	select {
	case <-p.ctx.Done():
		log.Printf("[worker %d] stopped", id)
		return
	case task, ok := <-p.tasks:
		if !ok {
			return
		}
		task()
	}
}

func (p *Pool) Submit(task Task) {
	select {
	case p.tasks <- task:
	default:
		log.Println("Worker pool is full, dropping task")
	}
}

// gracefull shutdown
func (p *Pool) Shutdown() {
	p.cancel()
	close(p.tasks)
	p.wg.Wait()
	log.Println("All workers stopped")
}
