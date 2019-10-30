package scraper

import (
	"net/http"
	"golang.org/x/net/html"
	"log"
)

/*
 the scraper looks through the contents of a website and creates a json for each type of
useful tag that it finds useful
for now : <div> , <para>, <pre>, <h[x]>, <a>
 */

func ScrapePage(webaddress string, chATags chan string, chTextData chan string) {

	resp, err := http.Get(webaddress)
	if err != nil {
		log.Printf("error reading in start page link: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	//log.Printf("WorkerId:%d started with webpage: %s", workerid, webaddress)

	tokenizer := html.NewTokenizer(resp.Body)

	chTexts := make(chan string, 1000)
	go readTokens(tokenizer, chATags, chTexts)

	mastertext:= ""
	for txt := range chTexts {
		mastertext += txt
	}

	chTextData<- mastertext
}

// for a given tokenizer, return back all the texts that you find inside the page
func readTokens(tokenizer *html.Tokenizer, chATags, chTexts chan string) {

	defer close(chTexts)

	var previousToken html.Token

	for {
		ttoken := tokenizer.Next()
		switch(ttoken) {
		case html.ErrorToken:
			log.Printf("error token hit")
			return

		case html.StartTagToken:
			previousToken = tokenizer.Token()
			tagnamebytes, taghasattr := tokenizer.TagName()
			if len(tagnamebytes) == 1 && tagnamebytes[0] == 'a' {
				// look for href attribute

				var tagkey, tagvalue []byte
				for (taghasattr) {

					tagkey, tagvalue, taghasattr = tokenizer.TagAttr()
					if string(tagkey) == "href" {
						chATags <- string(tagvalue)
						break
					}
				}
			}

		case html.EndTagToken:
			chTexts<- string('\n')

		case html.TextToken:
			if previousToken.Data == "script" {
				continue
			}

			chTexts<- html.UnescapeString(string(tokenizer.Text()))
		}
	}
}

