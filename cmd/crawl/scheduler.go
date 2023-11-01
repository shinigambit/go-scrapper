package crawl

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

type Scheduler struct {
	extractor         LinkExtractor
	poolSize          int
	channelBufferSize int
}

func NewScheduler(extractor LinkExtractor, poolSize, channelBufferSize int) Scheduler {
	if channelBufferSize < 1 {
		channelBufferSize = 1
	}
	if poolSize < 1 {
		poolSize = 1
	}
	return Scheduler{
		extractor:         extractor,
		poolSize:          poolSize,
		channelBufferSize: channelBufferSize,
	}
}

type Message struct {
	Url string
	Err error
}

func (s *Scheduler) ScheduleCrawl(ctx context.Context) (<-chan Message, error) {
	outChannel := make(chan Message)
	if s.poolSize < 1 {
		close(outChannel)
		return outChannel, fmt.Errorf("at least one worker is required to execute")
	}
	if s.channelBufferSize < 1 {
		close(outChannel)
		return outChannel, fmt.Errorf("a buffer size of at least one is required to execute")
	}
	effectiveDomain := s.extractor.EffectiveDomain()
	urlTracker := map[string]struct{}{
		effectiveDomain: {},
	}
	urlsChannel := make(chan string, s.channelBufferSize)
	urlsChannel <- effectiveDomain

	var once sync.Once
	var urlTrackerLock sync.RWMutex
	var pending atomic.Int64
	pending.Add(int64(len(urlsChannel)))
	for i := 0; i < s.poolSize; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					once.Do(func() {
						close(outChannel)
					})
					return
				case url, open := <-urlsChannel:
					if !open {
						return
					}
					links, err := s.extractor.Extract(url)
					if err != nil {
						outChannel <- Message{Err: err}
						if pending.Add(-1) == 0 {
							once.Do(func() {
								close(urlsChannel)
								close(outChannel)
							})
						}
						continue
					}
					outChannel <- Message{Url: url}

					go func() {
						urlTrackerLock.Lock()
						defer urlTrackerLock.Unlock()
						for _, link := range links {
							_, visited := urlTracker[link]
							if !visited {
								pending.Add(1)
								urlTracker[link] = struct{}{}
								urlsChannel <- link
							}
						}
						if pending.Add(-1) == 0 {
							once.Do(func() {
								close(urlsChannel)
								close(outChannel)
							})
						}
					}()
				}
			}
		}()
	}
	return outChannel, nil
}
