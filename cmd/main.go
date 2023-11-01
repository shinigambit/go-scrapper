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

	var urlGetter crawl.URLGetter
	throttling := os.Getenv("REQUEST_THROTTLING")
	if throttling == "false" {
		urlGetter = crawl.NewParallerRequestGetter(time.Second * 5)
	} else {
		urlGetter = crawl.NewHttpGetter(time.Second * 5)
	}
	extractor := crawl.NewUrlLinkExtractor(baseUrl, urlGetter)

	scheduler := crawl.NewScheduler(extractor, poolSize, 20)
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
