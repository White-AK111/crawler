package crawlerer

import (
	"context"
)

//Crawler - интерфейс (контракт) краулера
type Crawler interface {
	Scan(ctx context.Context, url string, parentUrl string, depth uint)
	ChanResult() <-chan CrawlResult
}

type Requester interface {
	Get(ctx context.Context, url string) (Page, error)
}

type Page interface {
	GetTitle() string
	GetLinks() []string
}

type CrawlResult struct {
	Err   error
	Title string
	Url   string
}
