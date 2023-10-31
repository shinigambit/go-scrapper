package crawl

import (
	"regexp"
	"strings"

	"crawler.parser.com/cmd/parse"
)

var urlStart *regexp.Regexp
var anchorSuffix *regexp.Regexp

func init() {
	urlStart = regexp.MustCompile("^https?://")
	anchorSuffix = regexp.MustCompile("/?#[^/]*$")
}

type LinkExtractor interface {
	EffectiveDomain() string
	Extract(string) ([]string, error)
}

type UrlLinkExtractor struct {
	domain    string
	urlGetter URLGetter
}

func NewUrlLinkExtractor(domain string, urlGetter URLGetter) LinkExtractor {
	if !strings.HasSuffix(domain, "/") {
		domain = domain + "/"
	}
	return &UrlLinkExtractor{
		domain:    domain,
		urlGetter: urlGetter,
	}
}

func (c *UrlLinkExtractor) EffectiveDomain() string {
	return c.domain
}

func (c *UrlLinkExtractor) Extract(url string) (links []string, err error) {
	response, err := c.urlGetter.Get(url)
	if err != nil {
		return
	}
	defer response.Close()

	rawLinks, err := parse.PageLinks(response)
	for _, s := range rawLinks {
		if link := c.toFullUrl(url, s); strings.HasPrefix(link, c.domain) {
			links = append(links, link)
		}
	}
	return
}

func (c *UrlLinkExtractor) toFullUrl(url, original string) string {
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
