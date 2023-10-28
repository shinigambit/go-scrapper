package main

import (
	"log"
	"os"
	"runtime"
	"sync"

	"crawler.parser.com/src/crawl"
)

func main() {
	poolSize := 4 * runtime.NumCPU()
	baseUrl := os.Getenv("CRAWLER_DOMAIN")
	if baseUrl == "" {
		baseUrl = "https://realpython.github.io/fake-jobs/"
	}

	// baseUrl = "https://parserdigital.com/"
	// baseUrl = "https://quotes.toscrape.com/"
	// baseUrl = "https://books.toscrape.com/"
	// baseUrl = "https://realpython.github.io/fake-jobs/"
	// baseUrl = "https://www.scrapethissite.com/"

	crawler := crawl.NewClient(baseUrl)
	urlsChannel := make(chan string, 50)
	effectiveUrl := crawler.EffectiveDomain()
	urlsChannel <- effectiveUrl
	urlTracker := map[string]struct{}{
		effectiveUrl: {},
	}

	runWorkers(crawler, poolSize, urlsChannel, urlTracker)
}

func runWorkers(crawler crawl.Client, poolSize int, urlsChannel chan string, urlTracker map[string]struct{}) {
	var wg sync.WaitGroup
	var urlTrackerLock sync.RWMutex
	wg.Add(len(urlsChannel))
	l := log.New(os.Stdout, "", 0)
	for i := 0; i < poolSize; i++ {
		go func() {
			for {
				url := <-urlsChannel
				l.Println(url)
				links, err := crawler.Request(url)
				if err != nil {
					log.New(os.Stderr, "error", log.Flags()).Printf("problem processing page (%v) links: %v", url, err)
					wg.Done()
					return
				}
				go func() {
					urlTrackerLock.Lock()
					defer urlTrackerLock.Unlock()
					for _, link := range links {
						_, visited := urlTracker[link]
						if !visited {
							wg.Add(1) // add one for every url that needs to be visited
							urlTracker[link] = struct{}{}
							urlsChannel <- link
						}
					}
					wg.Done() // url counts as visited when all its links have been sent to queue
				}()
			}
		}()
	}
	wg.Wait() // resume when all links have been visited
}
