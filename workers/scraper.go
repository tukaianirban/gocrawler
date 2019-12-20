package workers

import (
	"golang.org/x/net/html"
	"regexp"
	"prooftestideas/gocrawler/urlcache"
	"prooftestideas/gocrawler/page"
	"log"
)

/*
 the scraper looks through the contents of a website and creates a json for each type of
useful tag that it finds useful
for now : <div> , <para>, <pre>, <h[x]>, <a>
*/

// for a given tokenizer, return back all the texts that you find inside the page
func readTokens(tokenizer *html.Tokenizer, originalwebaddress string) *page.Page {

	var previousToken html.Token
	spaceregexp := regexp.MustCompile(`\s+`)

	thispage := &page.Page{
		WebAddress:		originalwebaddress,
		TagsA:			make([]page.TagA, 0),
		TagsMeta:		make([]page.TagMeta, 0),
		Data:			"",
	}

	for {
		ttoken := tokenizer.Next()

		if ttoken == html.ErrorToken {
			return thispage
		}

		tagnamebytes, _ := tokenizer.TagName()
		tagname := string(tagnamebytes)

		switch ttoken {

		case html.StartTagToken:
			previousToken = tokenizer.Token()

			if tagname == "a" {
				thispage.TagsA = append(thispage.TagsA, GetATags(tokenizer, originalwebaddress)...)
			}

			if tagname == "meta" {
				thispage.TagsMeta = append(thispage.TagsMeta , GetMetaTags(tokenizer)...)
			}

		case html.EndTagToken:
			if tagname == "div" || tagname == "p" || tagname == "pre" {
				thispage.Data = thispage.Data + string('\n')
			}

		case html.TextToken:
			if previousToken.Data == "script" {
				continue
			}

			textstring := html.UnescapeString(string(tokenizer.Text()))
			thispage.Data = thispage.Data + spaceregexp.ReplaceAllString(textstring, " ")
		}
	}
}

func GetATags(tokenizer *html.Tokenizer, originalurl string) []page.TagA {

	var tagkey, tagvalue []byte
	taghasattr := true

	atagslist := make([]page.TagA, 0)
	for taghasattr {

		tagkey, tagvalue, taghasattr = tokenizer.TagAttr()
		if string(tagkey) == "" {
			continue
		}

		if string(tagkey) == "href" {
			hrefpathtransform := validateAndTransform(originalurl, string(tagvalue))
			log.Printf("found new url: %s", hrefpathtransform)

			// add the newly-found weblink in the pages cache
			urlcache.AddDiscoveredWeblink(hrefpathtransform)

			atagslist = append(atagslist, map[string]string{string(tagkey): hrefpathtransform})

		}else {
			// add in the key-value for each tag-key : tag-value attribute
			atagslist = append(atagslist, map[string]string{string(tagkey): string(tagvalue)})
		}
	}

	return atagslist
}

func GetMetaTags(tokenizer *html.Tokenizer) []page.TagMeta {

	var tagkey, tagvalue []byte
	taghasattr := true

	atagslist := make([]page.TagMeta, 0)
	for taghasattr {

		tagkey, tagvalue, taghasattr = tokenizer.TagAttr()
		if string(tagkey) == "" {
			continue
		}

		// add in the key-value for each tag-key : tag-value attribute
		atagslist = append(atagslist, map[string]string{string(tagkey): string(tagvalue)})

		if string(tagkey) == "href" {
			// add the newly-found weblink in the pages cache
			urlcache.AddDiscoveredWeblink(string(tagvalue))
		}
	}

	return atagslist
}