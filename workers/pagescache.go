package workers

// master weblinks channel from where worker routines will pull and work on
var chMasterWebLinks = make(chan string, 10000)

var countDropped = 0

// todo: we need a better way to cache the new-found weblinks
// for now, we dont record in new found pages if the channel is at >0.75 of capacity
func AddDiscoveredWeblink(weblink string) {

	if len(chMasterWebLinks) < int ( 75 * cap(chMasterWebLinks) / 100) {
		chMasterWebLinks <- weblink
	}else {
		countDropped++
	}
}

func GetNewWeblink() string {

	return <-chMasterWebLinks
}

func GetPagesCacheSize() int {
	return len(chMasterWebLinks)
}

func GetPagesDroppedCount() int {
	return countDropped
}