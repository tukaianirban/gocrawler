package workers

import (
	"net/url"
)



/*
 * needed to take care of relative URls in the href attr of 'a' tag
 * todo: some non-english websites come up with unknown chars in the webaddress
 */
func validateAndTransform(baseaddr, reladdr string) string {

	relurl, err := url.Parse(reladdr)
	if err == nil && (relurl.Scheme == "http" || relurl.Scheme == "https") {
		return reladdr
	}

	baseurl, err := url.Parse(baseaddr)
	if err != nil {
		//
		// transformation failed since the base url address could not be resolved to an URL
		//
		return baseaddr
	}
	return baseurl.ResolveReference(relurl).String()
}
