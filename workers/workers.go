package workers

import (
	"time"
	"sync"
	"net/http"
	"prooftestideas/gocrawler/perf"
	"golang.org/x/net/html"
	"log"
)

type WorkerPool struct {
	PoolId      int
	State       bool
	tokenCount  int
	tokenLock   sync.RWMutex
	NextUrlChan chan string
}

func NewWorkerPool(poolId, max_no_workers int) *WorkerPool {

	log.Printf("New worker Pool created with id:%d and workerCount=%d", poolId, max_no_workers)

	return &WorkerPool{
		PoolId:      poolId,
		State:       true,
		tokenCount:  max_no_workers,
		NextUrlChan: make(chan string, 1000),
	}
}

func (self *WorkerPool) GetActiveWorkers() int {

	return self.tokenCount
}

func (self *WorkerPool) IsCapacityAvailable() bool {

	return self.tokenCount > 0
}

func (self *WorkerPool) StartWorkerPool() {

	log.Printf("Started running new worker pool id:%d", self.PoolId)

	for newlink := range self.NextUrlChan {

		err := self.scheduleWebLinkToWorker(newlink)
		if err != nil {

			// out of capacity to schedule the weblink to a worker
			self.NextUrlChan <- newlink
		}

	}
}

func (self *WorkerPool) scheduleWebLinkToWorker(newlink string) error {

	for retrydelay := 1; retrydelay < 10; retrydelay++ {

		workerid := self.GetWorkerToken()
		if workerid >= 0 {
			go self.ScrapePage(workerid, newlink)

			return nil
		}

		time.Sleep(time.Duration(retrydelay) * time.Second)
	}
	return ErrOutOfCapacity
}

func (self *WorkerPool) GetWorkerToken() int {

	if self.tokenCount == 0 {
		return -1
	}

	self.tokenLock.Lock()
	defer self.tokenLock.Unlock()

	workerId := self.tokenCount

	self.tokenCount--

	return workerId
}

func (self *WorkerPool)ReturnWorkerToken() {

	self.tokenLock.Lock()
	defer self.tokenLock.Unlock()

	self.tokenCount++
}

func (self *WorkerPool)ScrapePage(workerid int, webaddress string) {

	defer self.ReturnWorkerToken()

	resp, err := http.Get(webaddress)
	if err != nil {
		log.Printf("error reading in start page link: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	//log.Printf("WorkerId:%d started with webpage: %s", workerid, webaddress)

	defer perf.AddPageIndexed()

	tokenizer := html.NewTokenizer(resp.Body)

	// todo: store this in a database / inline cache
	chTexts := make(chan string, 5000)
	go readTokens(tokenizer, chTexts)

	mastertext:= ""
	for txt := range chTexts {
		mastertext += txt
	}

	//log.Printf("workerId: %d weblink: %s textdump length: %d", workerid, webaddress, len(mastertext))
}