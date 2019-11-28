package epubln

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func CrawlAllByLabel(ctx context.Context, label string) (Downloader, error) {
	return CrawlAllInLabelPage(ctx, "https://epubln.blogspot.com/search/label/"+label)
}

func CrawlAllInLabelPage(ctx context.Context, rawurl string) (Downloader, error) {
	doc, err := fetchPageByURL(ctx, rawurl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page content: %v", err)
	}

	var wg sync.WaitGroup
	dg := newGroupDownloader()
	doc.Find("a.anes").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("href")
		if !ok {
			return
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if job, err := CrawlFromPostPage(ctx, url); err == nil {
				dg.Add(job)
			}
		}()
	})
	wg.Wait()
	return dg, nil
}

func CrawlFromPostPage(ctx context.Context, rawurl string) (Downloader, error) {
	doc, err := fetchPageByURL(ctx, rawurl)
	if err != nil {
		log.Printf("Failed to read from post page %v: %v\n", rawurl, err)
		return nil, err
	}

	dg := newGroupDownloader()
	doc.Find(".cover a").Each(func(i int, s *goquery.Selection) {
		url, err := url.Parse(s.AttrOr("href", ""))
		if err != nil {
			log.Printf("Failed to parse download String in post page %v: %v\n", rawurl, err)
			return
		}

		d, err := newDownloader(url)
		if err != nil {
			return
		}
		dg.Add(d)
	})
	return dg, nil
}

func fetchPageByURL(ctx context.Context, rawurl string) (*goquery.Document, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawurl, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create new HTTP request: %v", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot do HTTP request: %v", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP body: %v", err)
	}
	return doc, nil
}
