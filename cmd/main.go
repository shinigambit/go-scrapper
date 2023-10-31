package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"crawler.parser.com/cmd/crawl"
)

func main() {
	poolSize := 4 * runtime.NumCPU()
	baseUrl := os.Getenv("CRAWLER_DOMAIN")
	if baseUrl == "" {
		baseUrl = "https://www.scrapethissite.com/"
	}

	// baseUrl = "https://parserdigital.com/"
	// baseUrl = "https://quotes.toscrape.com/"
	// baseUrl = "https://books.toscrape.com/"
	// baseUrl = "https://realpython.github.io/fake-jobs/"
	// baseUrl = "https://www.scrapethissite.com/"
	// baseUrl = "https://parserdigital.com/contact-us/javascript:void(0)"

	urlGetter := crawl.NewHttpGetter(time.Second * 3)
	extractor := crawl.NewUrlLinkExtractor(baseUrl, urlGetter)

	scheduler := crawl.NewScheduler(extractor, poolSize, 50)
	resultsChannel, err := scheduler.ScheduleCrawl(context.Background())
	if err != nil {
		fmt.Println("problem initiating:", err)
		return
	}
	for result := range resultsChannel {
		if result.Url != "" {
			fmt.Println(result.Url)
		}
		if result.Err != nil {
			fmt.Println(result.Err)
		}
	}
}
