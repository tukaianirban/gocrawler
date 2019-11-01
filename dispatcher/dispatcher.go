package dispatcher

import (
	"sync"
	"prooftestideas/gocrawler/workers"
	"time"
	"prooftestideas/gocrawler/urlcache"
	"runtime"
	"net/url"
	"prooftestideas/gocrawler/perf"
	"log"
)

// a dispatcher schedules the next URL to be loaded to a worker pool
// it pulls the next URL from the urlcache and tries to schedule it for scraping
// it takes into account various factors like frequency of hits on a certain domain (todo), availability of worker pools, etc
// before dispatching the next URL to a worker pool

type Dispatcher struct {
	poolsmap      map[int]*workers.WorkerPool
	poolLock      sync.RWMutex
	maxPoolsCount int
	chDone        chan bool
}

func NewDispatcher(maxPoolsCount int) *Dispatcher {

	return &Dispatcher{
		poolsmap:      make(map[int]*workers.WorkerPool),
		maxPoolsCount: maxPoolsCount,
		chDone:        make(chan bool),
	}
}

// we dont want many packages to call the urlcache package, hence this function
func (self *Dispatcher) StartDispatcher(starturl string, chDone chan bool, maxPoolsCount int) {

	urlcache.PutBackWeblink(starturl)
	self.dispatcher()

}

// the stop dispatcher will signal all pools to stop.
// then it will monitor until all pools are dead (deleted from the map by their shutdown routine)
// this function will return only when all the pools are dead
func (self *Dispatcher) StopDispatcher() {

	for poolid, _ := range self.poolsmap {
		self.stopWorkerPool(poolid)
	}

	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {

		if len(self.poolsmap) == 0 {
			return
		}
	}
}

func (self *Dispatcher) dispatcher() {

	isAlive := true

	for isAlive {

		// pull the next URL to load, from the urlcache
		select {
		case nextURL := <-urlcache.GetNextWeblink():

			if !validateURL(nextURL) {
				// mark the URL as invalid, count it, and drop it
				perf.AddPageInvalidWeblink()
				break
			}

			isScheduled := self.scheduleURLToPool(nextURL)
			if isScheduled {
				break
			}

			// do some throttling here


		case <-self.chDone:
			self.StopDispatcher()
			isAlive = false
		}
	}

}

// start a new worker pool
func (self *Dispatcher) startNewWorkerPool(nextURL string) {

	currentLen := len(self.poolsmap)

	max_no_workers := runtime.NumCPU()

	self.poolLock.Lock()
	defer self.poolLock.Unlock()

	self.poolsmap[currentLen] = workers.NewWorkerPool(currentLen, max_no_workers)
	go self.poolsmap[currentLen].StartWorkerPool()

	if nextURL != "" {
		self.poolsmap[currentLen].NextUrlChan <- nextURL
	}
}

func (self *Dispatcher) scheduleURLToPool(url string) bool {

	for id := range self.poolsmap {

		if self.poolsmap[id].State && self.poolsmap[id].IsCapacityAvailable() {

			self.poolsmap[id].NextUrlChan <- url
			return true
		}
	}

	//log.Printf("Failed to schedule url to any worker pool")

	// see if a new pool can be started
	if len(self.poolsmap) < self.maxPoolsCount {

		log.Printf("Starting new pool with id:%d", len(self.poolsmap))
		self.startNewWorkerPool(url)

		return true
	}

	// no pools available
	// Put the weblink back to urlcache
	urlcache.PutBackWeblink(url)

	return false
}

// gracefully shutdown a worker pool
func (self *Dispatcher) stopWorkerPool(poolid int) {

	self.poolLock.Lock()
	defer self.poolLock.Unlock()

	// set the pool to deactivated
	self.poolsmap[poolid].State = false
	close(self.poolsmap[poolid].NextUrlChan)

	// start a monitoring routine to remove the entry from the poolsmap when the workerpool has completely
	// finished its working
	go func(poolid int) {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {

			toDelete := false

			self.poolLock.RLock()
			if self.poolsmap[poolid].GetActiveWorkers() == 0 && !self.poolsmap[poolid].State {
				toDelete = true
			}
			self.poolLock.RUnlock()

			if toDelete {
				self.poolLock.Lock()
				delete(self.poolsmap, poolid)
				self.poolLock.Unlock()
			}
		}
	}(poolid)
}

func validateURL(rawurl string) bool {

	u, err := url.Parse(rawurl)
	if err != nil {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	if u.Host == "" {
		return false
	}

	return true
}