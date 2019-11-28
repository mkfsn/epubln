package main

import (
	"context"
	"log"
	"strings"

	"github.com/mkfsn/epubln"
)

func main() {
	//DownloadAllByLabel(`完結-龍孃七七七埋藏的寶藏`)
	//DownloadAllInLabelPage(`https://epubln.blogspot.com/search/label/完結-笨蛋，測驗，召喚獸`)
	//DownloadAllByLabel(`完結-虛軸少女`)
	//DownloadAllByLabel(`完結-黃昏色的詠使`)
	//DownloadAllByLabel(`待續-刀劍神域-刀劍神域`)

	downloader, err := epubln.CrawlAllByLabel(context.Background(), `完結-龍孃七七七埋藏的寶藏`)
	if err != nil {
		log.Fatalf("Failed to crawl all by label: %v\n", err)
	}

	for _, link := range strings.Split(downloader.String(), " ") {
		log.Println(link)
	}

	if err := downloader.Download(); err != nil {
		log.Fatalf("Failed to download: %v\n", err)
	}
}
