package workers

import (
	"testing"
	"runtime"
	"prooftestideas/gocrawler/page"
	"prooftestideas/gocrawler/urlcache"
)

func TestWorkerPool_GetWorkerToken(t *testing.T) {

	maxWorkerCount := runtime.NumCPU()
	testPool := NewWorkerPool(1, maxWorkerCount)

	for i:=0; i<maxWorkerCount; i++ {
		t.Logf("received worker token: %d", testPool.GetWorkerToken())
	}
	t.Logf("all tokens drained from pool")

	t.Logf("next token request: %d", testPool.GetWorkerToken())
	t.Logf("returning more tokens than obtained ...")
	for i:=0; i<maxWorkerCount + 2; i++ {
		testPool.ReturnWorkerToken()
	}

	t.Logf("drain all tokens again...")
	for i:=0; i<maxWorkerCount; i++ {
		t.Logf("received worker token: %d", testPool.GetWorkerToken())
	}
	t.Logf("all tokens drained from pool")
}

func TestWorkerPool_ScrapePage(t *testing.T) {

	workerId := 1
	testPageURL := "https://en.wikipedia.org/wiki/Amazon_Web_Services"

	testPool := NewWorkerPool(1, runtime.NumCPU())
	if err := urlcache.InitCache(); err!= nil {
		t.Fatalf("error initiating the pages redis cache: %s", err.Error())
	}

	dummyPagesCacheRegister := func(urladdress string, pageScraped *page.Page) {

		t.Logf("scraped a new page:")
		t.Logf("url: %s", urladdress)
		t.Logf("Page has %d Meta tags", len(pageScraped.TagsMeta))
		t.Logf("Page has %d A tags", len(pageScraped.TagsA))
	}

	testPool.ScrapePage(workerId, testPageURL, dummyPagesCacheRegister)

}