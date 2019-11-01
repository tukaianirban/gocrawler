package workers

import (
	"log"
	"net/http"
	"golang.org/x/net/html"
	"prooftestideas/gocrawler/perf"
	"prooftestideas/gocrawler/urlcache"
)

/*
 * given a page's web address, search through the page for <a> start tokens and feed them into the queue
 * for next searching worker
 */
// todo: deprecated
func SearchThroughPage(workerid int, webaddress string) {

	// signal that this worker is done
	//defer ReturnWorkerToken()

	// add a count for having indexed this page
	defer perf.AddPageIndexed()

	resp, err := http.Get(webaddress)
	if err != nil {
		log.Printf("error reading in start page link: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	//log.Printf("WorkerId:%d started with webpage: %s", workerid, webaddress)

	tokenizer := html.NewTokenizer(resp.Body)

	chWebLinks := make(chan string, 1000)
	go lookForATags(tokenizer, chWebLinks)

	//
	// we could merge all the found weblinks directly into the master channel for weblinks
	// but this is what we are doing for now.
	// we will use this place for some sanity checks on the weblinks found in this page
	// todo: make a better implementation here

	count := 0
	for weblink := range chWebLinks {

		urlcache.AddDiscoveredWeblink(weblink)
		count++
	}

	//log.Printf("WorkerId: %d: Found %d references from page: %s", workerid, count, webaddress)
}

func lookForATags(tokenizer *html.Tokenizer, chWebLinks chan string) {

	defer close(chWebLinks)

	for {
		ttoken := tokenizer.Next()
		if ttoken == html.ErrorToken {
			return
		}

		// look for start tags only
		if ttoken != html.StartTagToken {
			continue
		}

		tagnamebytes, taghasattr := tokenizer.TagName()
		if len(tagnamebytes) == 1 && tagnamebytes[0] == 'a' {
			// look for href attribute

			var tagkey, tagvalue []byte
			for (taghasattr) {

				tagkey, tagvalue, taghasattr = tokenizer.TagAttr()
				if string(tagkey) == "href" {
					chWebLinks <- string(tagvalue)
					break
				}
			}
		}
	}
}
