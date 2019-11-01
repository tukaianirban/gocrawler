package pagescache

// inline cache of the next URLs to crawl through
// the cache contains a unique list of websites and the number of times they were hit
var chMasterWebLinks = make(chan string, 10000)

var countDropped = 0

// todo: we need a better way to cache the new-found weblinks
// for now, we dont record in new found pages if the channel is at >0.75 of capacity
func AddDiscoveredWeblink(weblink string) {

	if FoundNewIndexedPage(weblink) {
		return
	}

	// for now, the inline cache of nextURLToLoad is maxed at 75% of capacity
	if len(chMasterWebLinks) < int ( 75 * cap(chMasterWebLinks) / 100) {
		chMasterWebLinks <- weblink
	}else {
		countDropped++
	}
}

// unconditionally put back the weblink into the cache
// ex situation: when scheduler fails to schedule the weblink to a worker pool
func PutBackWeblink(weblink string) {
	chMasterWebLinks <- weblink
}

func GetNextWeblink() chan string {

	return chMasterWebLinks
}

func GetPagesCacheSize() int {
	return len(chMasterWebLinks)
}

func GetPagesDroppedCount() int {
	return countDropped
}