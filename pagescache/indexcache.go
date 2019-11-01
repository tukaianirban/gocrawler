package pagescache

import "sync"

// a local store of indexed pages and the number of times they were hit

var indexedpages = make(map[string]int)
var indexLock sync.RWMutex

// when a new page is found, update the hit count of the page
// returns a bool of whether the page is already known or not
func FoundNewIndexedPage(uri string) bool {

	indexLock.Lock()
	defer indexLock.Unlock()

	existinghits := indexedpages[uri] > 0

	indexedpages[uri]++

	return existinghits
}
