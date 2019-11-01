package main

import (
	"log"
	"prooftestideas/gocrawler/perf"
	"time"
	"flag"
	"prooftestideas/gocrawler/pagescache"
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

	chDone := make(chan bool, 10)
	dispatcher1 := dispatcher.NewDispatcher(1)
	go dispatcher1.StartDispatcher(*startPageLink, chDone, 1)

	go PrintPerformanceStats()

	// essentially, for now, this will run forever
	<-chDone
}

func PrintPerformanceStats() {

	for {
		log.Printf("Pages indexed = %d PagesCache size = %d PagesCache dropped = %d InvalidPages = %d",
			perf.GetPagesIndexed(), pagescache.GetPagesCacheSize(), pagescache.GetPagesDroppedCount(), perf.GetPageInvalidWeblinkCount())

		time.Sleep(10 * time.Second)
	}
}


