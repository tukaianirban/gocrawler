package main

import (
	"log"
	"prooftestideas/gocrawler/perf"
	"time"
	"flag"
	"prooftestideas/gocrawler/urlcache"
	"prooftestideas/gocrawler/dispatcher"
)

var startPageLink *string

// randomly chosen page to start the crawler from
var DEFAULT_START_PAGE = "https://en.wikipedia.org/wiki/Amazon_Web_Services"

func init() {

	startPageLink = flag.String("startpage", DEFAULT_START_PAGE, "page to start the crawler from")

	flag.Parse()

}

func main() {

	log.Println("Starting ...")

	if err := urlcache.InitCache(); err!=nil {
		log.Fatalf("error setting up the url cache:%s", err.Error())
	}

	chDone := make(chan bool, 10)
	dispatcher1 := dispatcher.NewDispatcher(1)
	go dispatcher1.StartDispatcher(*startPageLink, chDone, 1)

	go PrintPerformanceStats()

	chFinishPageCacheDump := make(chan struct{})
	go dumpWebPagesInCache(chFinishPageCacheDump)

	// essentially, for now, this will run forever
	<-chDone

	chFinishPageCacheDump<- struct{}{}
}

func PrintPerformanceStats() {

	for {
		log.Printf("Pages indexed = %d PagesCache size = %d PagesCache dropped = %d InvalidPages = %d urlCache size=%d",
			perf.GetPagesIndexed(), urlcache.GetPagesCacheSize(), urlcache.GetPagesDroppedCount(), perf.GetPageInvalidWeblinkCount(), urlcache.GetUrlCacheSize())

		time.Sleep(10 * time.Second)
	}
}

func dumpWebPagesInCache(chDone chan struct{}) {

	for {
		select {
			case <-chDone:
				return

			case <-time.After(1 * time.Second):

		}

		nextPage := urlcache.GetNextPage()
		if nextPage == nil {
			continue
		}

		log.Printf("next page scraped: %s", nextPage.WebAddress)
	}
}
