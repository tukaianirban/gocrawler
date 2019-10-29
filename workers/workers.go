package workers

import (
	"sync"
	"errors"
	"log"
	"time"
	"net/url"
	"prooftestideas/gocrawler/perf"
)



var ErrOutOfCapacity = errors.New("out of worker capacity")
var ErrInvalidWeblink = errors.New("invalid weblink")

func StartWorkerPool(wg *sync.WaitGroup) {

	for newlink := range chMasterWebLinks {

		err := scheduleWebLinkToWorker(newlink)
		switch err {
			case nil:

			case ErrInvalidWeblink:
				//log.Printf("invalid weblink. skip it")

			case ErrOutOfCapacity:
				log.Printf("timed out waiting for pushing into the worker pool. Put the weblink back in there")
				chMasterWebLinks <- newlink
				log.Printf(err.Error())

		}
	}

	wg.Done()
}

func scheduleWebLinkToWorker(newlink string) error {

	// sanity check the weblink
	if !validateURL(newlink) {
		perf.AddPageInvalidWeblink()
		return ErrInvalidWeblink
	}

	for retrydelay:=1; retrydelay < 10; retrydelay++ {

		token := GetWorkerToken()
		if token>=0 {
			go SearchThroughPage(token, newlink)

			return nil
		}

		time.Sleep(time.Duration(retrydelay) * time.Second)
	}
	return ErrOutOfCapacity
}

func validateURL(rawurl string) bool {

	u, err := url.Parse(rawurl)
	if err!=nil {
		return false
	}

	if u.Scheme!="http" && u.Scheme!="https" {
		return false
	}

	if u.Host=="" {
		return false
	}

	return true
}