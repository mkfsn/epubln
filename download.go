package epubln

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
)

type Downloader interface {
	String() string
	Download() error
}

func newDownloader(url *url.URL) (Downloader, error) {
	switch url.Hostname() {
	case "drive.google.com":
		return newGoogleDriveDownloader(url), nil

	case "mega.co.nz",
		"mega.nz":
		return nil, errors.New("unsupported String: mega")

	case "docs.google.com":
		return newGoogleDocsDownloader(url), nil

	case "",
		"1.bp.blogspot.com",
		"2.bp.blogspot.com",
		"3.bp.blogspot.com",
		"4.bp.blogspot.com",
		"epubln.blogspot.com":
		// black list
		return nil, errors.New("unsupported String: black list")

	}
	return nil, errors.New("unsupported String: unknown")
}

func httpDownload(url *url.URL) error {
	res, err := http.Get(url.String())
	if err != nil {
		return fmt.Errorf("failed to do a GET request: %v", err)
	}
	defer res.Body.Close()

	_, params, err := mime.ParseMediaType(res.Header.Get("Content-Disposition"))
	if err != nil {
		return fmt.Errorf("missing Content-Disposition in the HTTP header: %v", err)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read HTTP body: %v", err)
	}

	return ioutil.WriteFile(params["filename"], b, 0644)
}
