package urlcache

import "log"

// inline cache of the next URLs to crawl through
// the cache contains a unique list of websites and the number of times they were hit
var chMasterWebLinks = make(chan string, 10000)
var countDropped = 0

var redisClient *RedisClient

func InitCache() error {

	var err error

	redisClient, err = NewRedisCache("urlHits")

	return err
}

// todo: if the url exists in redis cache, then increment its hit count
// todo: 		- else, add it into chMasterWebLinks to schedule it for scraping
//
func AddDiscoveredWeblink(weblink string) {

	if redisClient == nil {
		log.Fatalf("cache is not yet initialized !")
	}

	isExisted, err := redisClient.SetURLHit(weblink)
	if err != nil {
		log.Printf("error setting hit count on existing url in cache:%s", err.Error())
	}

	if !isExisted {
		if len(chMasterWebLinks) < int(75*cap(chMasterWebLinks)/100) {
			chMasterWebLinks <- weblink
		} else {
			countDropped++
		}
	}

	// for now, the inline cache of nextURLToLoad is maxed at 75% of capacity

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

func GetUrlCacheSize() int64 {
	return redisClient.GetCacheSize()
}

func GetPagesDroppedCount() int {
	return countDropped
}

// this is a newly indexed page scraped up by a crawler
// todo: add the url to redis cache
// todo: add the url + data to database
func AddDiscoveredPage(url, data string) {

	if err := redisClient.StoreURL(url); err != nil {
		log.Printf("error storing URL:%s in cache: %s", url, err.Error())
	}
	//log.Printf("stored new url + page contents in cache")
}
