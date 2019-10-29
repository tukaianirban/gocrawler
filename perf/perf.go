package perf

import (
	"sync"
)

var countURLsCrawled int = 0
var countURLsInvalid int = 0

var perflock sync.RWMutex

func AddPageIndexed() {

	perflock.Lock()
	countURLsCrawled++
	perflock.Unlock()
}

func GetPagesIndexed() int {

	perflock.RLock()
	defer perflock.RUnlock()

	return countURLsCrawled

}

func AddPageInvalidWeblink() {

	perflock.Lock()
	countURLsInvalid++
	perflock.Unlock()
}

func GetPageInvalidWeblinkCount() int {

	perflock.RLock()
	defer perflock.RUnlock()

	return countURLsInvalid
}

