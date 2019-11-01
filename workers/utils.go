package workers

import "net/url"



/*
 * needed to take care of relative URls in the href attr of 'a' tag
 */
func validateAndTransform(originalurl, relurl string) string {

	relurltrans, err := url.Parse(relurl)
	if err == nil && (relurltrans.Scheme == "http" || relurltrans.Scheme == "https") {
		return relurl
	}

	// the relurl is (by here) not a valid url
	// obtain required stuff from originalurl
	originalurltrans, err := url.Parse(originalurl)
	if err != nil {
		return ""
	}

	originalurltrans.RawQuery = ""
	originalurltrans.Fragment = ""
	originalurltrans.Path = url.PathEscape(relurl)

	return originalurltrans.String()
}
