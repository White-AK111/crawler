package crawlerer

import (
	"context"
)

//Crawler - интерфейс (контракт) краулера
type Crawler interface {
	Scan(ctx context.Context, url string, parentUrl string, maxDepth *int64, depth int64)
	ChanResult() <-chan CrawlResult
	ToChanResult(CrawlResult)
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
	Info  string
	Title string
	Url   string
}
