package crawl

import (
	"net/http"
	"regexp"
	"strings"

	"crawler.parser.com/src/parse"
)

var urlStart *regexp.Regexp
var anchorSuffix *regexp.Regexp

func init() {
	urlStart = regexp.MustCompile("^https?://")
	anchorSuffix = regexp.MustCompile("/?#[^/]*$")
}

type Client struct {
	domain string
	getUrl func(string) (*http.Response, error)
}

func NewClient(domain string) Client {
	if !strings.HasSuffix(domain, "/") {
		domain = domain + "/"
	}
	return Client{
		domain: domain,
		getUrl: http.Get,
	}
}

func (c *Client) EffectiveDomain() string {
	return c.domain
}

func (c *Client) Request(url string) (links []string, err error) {
	response, err := c.getUrl(url)
	if err != nil {
		return
	}
	defer response.Body.Close()
	if err != nil {
		return
	}

	rawLinks, err := parse.PageLinks(response.Body)
	for _, s := range rawLinks {
		if link := c.toAbsoluteLink(url, s); strings.HasPrefix(link, c.domain) {
			links = append(links, link)
		}
	}
	return
}

func (c *Client) toAbsoluteLink(url, original string) string {
	// remove anchor from url
	original = anchorSuffix.ReplaceAllString(original, "")
	// absolute path
	if strings.HasPrefix(original, "/") {
		return c.domain + original[1:]
	}
	// full url
	if urlStart.MatchString(original) {
		return original
	}
	withoutLeadingSlash := strings.TrimPrefix(original, "/")
	// relative path
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	return url + withoutLeadingSlash
}
