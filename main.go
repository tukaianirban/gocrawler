package main

import (
	"log"
	"sync"
	"prooftestideas/gocrawler/workers"
	"prooftestideas/gocrawler/perf"
	"time"
	"flag"
	"runtime"
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

	// set the number of workers to be working
	cpucount := runtime.NumCPU()
	workers.SetMaxWorkersCount(cpucount)

	workers.AddDiscoveredWeblink(*startPageLink)

	var wg sync.WaitGroup

	go workers.StartWorkerPool(&wg)
	wg.Add(1)


	go PrintPerformanceStats()
	// essentially, for now, this will run forever
	wg.Wait()

}

func PrintPerformanceStats() {

	for {
		log.Printf("Pages indexed = %d PagesCache size = %d PagesCache dropped = %d InvalidPages = %d",
			perf.GetPagesIndexed(), workers.GetPagesCacheSize(), workers.GetPagesDroppedCount(), perf.GetPageInvalidWeblinkCount())

		time.Sleep(10 * time.Second)
	}
}


