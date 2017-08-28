package server

import (
	"net"
)

type WorkerPool struct {
	WorkerFunc func(c net.Conn)

	JobChan chan net.Conn
}

func (workerPool *WorkerPool) Start() {
	for i := 0; i <= 3000; i++ {
		go func(jobs <-chan net.Conn) {
			workerPool.WorkerFunc(<-jobs)
		}(workerPool.JobChan)
	}
}

func (workerPool *WorkerPool) Serve(connection net.Conn) {
	workerPool.JobChan <- connection
}
