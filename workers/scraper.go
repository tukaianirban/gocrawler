package workers

import (
	"golang.org/x/net/html"
	"regexp"
	"prooftestideas/gocrawler/pagescache"
)

/*
 the scraper looks through the contents of a website and creates a json for each type of
useful tag that it finds useful
for now : <div> , <para>, <pre>, <h[x]>, <a>
 */



// for a given tokenizer, return back all the texts that you find inside the page
func readTokens(tokenizer *html.Tokenizer, chTexts chan string) {

	defer close(chTexts)

	var previousToken html.Token
	spaceregexp := regexp.MustCompile(`\s+`)

	for {
		ttoken := tokenizer.Next()

		if ttoken == html.ErrorToken {
			return
		}

		tagnamebytes, taghasattr := tokenizer.TagName()
		tagname := string(tagnamebytes)

		switch(ttoken) {

		case html.StartTagToken:
			previousToken = tokenizer.Token()

			if tagname == "a" {

				var tagkey, tagvalue []byte
				for taghasattr {

					tagkey, tagvalue, taghasattr = tokenizer.TagAttr()
					if string(tagkey) == "href" {
						// add the newly-found weblink in the pages cache
						pagescache.AddDiscoveredWeblink(string(tagvalue))
						break
					}
				}
			}

		case html.EndTagToken:
			if tagname == "div" || tagname == "p" || tagname == "pre" {
				chTexts<- string('\n')
			}

		case html.TextToken:
			if previousToken.Data == "script" {
				continue
			}

			textstring := html.UnescapeString(string(tokenizer.Text()))
			chTexts<- spaceregexp.ReplaceAllString(textstring, " ")
		}
	}
}

