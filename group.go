package epubln

import (
	"log"
	"strings"
	"sync"
)

type groupDownloader struct {
	sync.WaitGroup
	jobs []Downloader
}

func newGroupDownloader() *groupDownloader {
	return &groupDownloader{}
}

func (g *groupDownloader) Add(job Downloader) {
	g.jobs = append(g.jobs, job)
}

func (g *groupDownloader) Download() error {
	g.WaitGroup.Add(len(g.jobs))
	for _, downloader := range g.jobs {
		downloader := downloader
		go func() {
			defer g.WaitGroup.Done()
			if err := downloader.Download(); err != nil {
				log.Printf("Failed to download %s: %v\n", downloader.String(), err)
			}
		}()
	}
	g.WaitGroup.Wait()
	return nil
}

func (g *groupDownloader) String() string {
	res := make([]string, 0, len(g.jobs))
	for _, job := range g.jobs {
		res = append(res, job.String())
	}
	return strings.Join(res, " ")
}
