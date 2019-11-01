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

	// essentially, for now, this will run forever
	<-chDone
}

func PrintPerformanceStats() {

	for {
		log.Printf("Pages indexed = %d PagesCache size = %d PagesCache dropped = %d InvalidPages = %d urlCache size=%d",
			perf.GetPagesIndexed(), urlcache.GetPagesCacheSize(), urlcache.GetPagesDroppedCount(), perf.GetPageInvalidWeblinkCount(), urlcache.GetUrlCacheSize())

		time.Sleep(10 * time.Second)
	}
}


