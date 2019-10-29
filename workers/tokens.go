package workers

import (
	"sync"
	"log"
)

var MAX_WORKERS = 10

var lock sync.RWMutex
var workerMaster = &TokenBucket{Tokencount:MAX_WORKERS}

type TokenBucket struct {
	Tokencount int
}

func GetWorkerToken() int {

	lock.Lock()
	defer lock.Unlock()

	if workerMaster.Tokencount == 0 {
		return -1
	}
	workerId:= workerMaster.Tokencount

	workerMaster.Tokencount--

	return workerId
}

func WorkerDone() {

	lock.Lock()

	if workerMaster.Tokencount<MAX_WORKERS {
		workerMaster.Tokencount++
	}

	lock.Unlock()
}

func SetMaxWorkersCount(count int) {

	MAX_WORKERS = count
	log.Printf("Resetting count of workers to %d", MAX_WORKERS)
}

func GetMaxWorkersCount() int {
	return MAX_WORKERS
}