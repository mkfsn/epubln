package epubln

import (
	"net/url"
	"strings"
)

type googleDriveDownloader struct {
	url *url.URL
}

func newGoogleDriveDownloader(u *url.URL) *googleDriveDownloader {
	if u.Path == "/uc" {
		return &googleDriveDownloader{u}
	}

	// TODO: how to ensure the path is in this format: /file/d/{id}/edit ?
	id := strings.TrimPrefix(u.Path, "/file/d/")
	id = strings.TrimSuffix(id, "/edit")
	id = strings.TrimSuffix(id, "/view")

	q := make(url.Values)
	q.Add("export", "download")
	q.Add("id", id)
	u.Path = "/uc"
	u.RawQuery = q.Encode()

	return &googleDriveDownloader{u}
}

func (g *googleDriveDownloader) Download() error {
	return httpDownload(g.url)
}

func (g *googleDriveDownloader) String() string {
	return g.url.String()
}

type googleDocsDownloader struct {
	url *url.URL
}

func newGoogleDocsDownloader(u *url.URL) *googleDocsDownloader {
	// TODO: how to ensure the path is in this format: /file/d/{id}/edit ?
	id := strings.TrimPrefix(u.Path, "/file/d/")
	id = strings.TrimSuffix(id, "/edit")

	q := make(url.Values)
	q.Add("export", "download")
	q.Add("id", id)

	u.Host = "drive.google.com"
	u.Path = "/uc"
	u.RawQuery = q.Encode()

	return &googleDocsDownloader{u}
}

func (g *googleDocsDownloader) Download() error {
	return httpDownload(g.url)
}

func (g *googleDocsDownloader) String() string {
	return g.url.String()
}
