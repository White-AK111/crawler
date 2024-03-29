package crawler

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/t0pep0/GB_best_go1/crawlerer"
)

// crawler struct
type crawler struct {
	r       crawlerer.Requester
	res     chan crawlerer.CrawlResult
	visited map[string]struct{}
	mu      sync.RWMutex
	wg      *sync.WaitGroup
}

// NewCrawler func for initialize new crawler struct
func NewCrawler(r crawlerer.Requester) *crawler {
	var wg sync.WaitGroup
	wg.Add(1)

	crawler := &crawler{
		r:       r,
		res:     make(chan crawlerer.CrawlResult),
		visited: make(map[string]struct{}),
		mu:      sync.RWMutex{},
		wg:      &wg,
	}

	go func() {
		wg.Wait()
		crawler.ToChanResult(crawlerer.CrawlResult{
			Info: "All URL's scanned",
		})
	}()

	return crawler
}

func (c *crawler) ToChanResult(crawResult crawlerer.CrawlResult) {
	c.res <- crawResult
}

// Scan method for crawler
func (c *crawler) Scan(ctx context.Context, url string, parentUrl string, maxDepth *int64, depth int64) {
	defer c.wg.Done()

	if depth > *maxDepth { //Проверяем то, что есть запас по глубине
		return
	}

	// Crutch for have a little more live links
	if !strings.HasPrefix(url, "http") {
		lInd := strings.LastIndex(parentUrl, "/")
		url = parentUrl[:lInd+1] + url
	}

	c.mu.RLock()
	_, ok := c.visited[url] //Проверяем, что мы ещё не смотрели эту страницу
	c.mu.RUnlock()
	if ok {
		return
	}
	select {
	case <-ctx.Done(): //Если контекст завершен - прекращаем выполнение
		return
	default:
		page, err := c.r.Get(ctx, url) //Запрашиваем страницу через Requester
		if err != nil {
			c.res <- crawlerer.CrawlResult{Err: err} //Записываем ошибку в канал
			return
		}
		c.mu.Lock()
		c.visited[url] = struct{}{} //Помечаем страницу просмотренной
		c.mu.Unlock()
		c.res <- crawlerer.CrawlResult{ //Отправляем результаты в канал
			Title: page.GetTitle(),
			Url:   url,
		}

		for _, link := range page.GetLinks() {
			c.wg.Add(1)
			go c.Scan(ctx, link, url, maxDepth, depth+1) //На все полученные ссылки запускаем новую рутину сборки
		}
	}
}

// ChanResult method for put result in channel
func (c *crawler) ChanResult() <-chan crawlerer.CrawlResult {
	return c.res
}

type requester struct {
	timeout time.Duration
}

func NewRequester(timeout time.Duration) requester {
	return requester{timeout: timeout}
}

func (r requester) Get(ctx context.Context, url string) (crawlerer.Page, error) {
	select {
	case <-ctx.Done():
		return nil, nil
	default:
		cl := &http.Client{
			Timeout: r.timeout,
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		body, err := cl.Do(req)
		if err != nil {
			return nil, err
		}
		defer body.Body.Close()
		page, err := NewPage(body.Body)
		if err != nil {
			return nil, err
		}
		return page, nil
	}

	// govet
	// delete unreachable return
	// unreachable: unreachable code (govet)
}

type page struct {
	doc *goquery.Document
}

func NewPage(raw io.Reader) (crawlerer.Page, error) {
	doc, err := goquery.NewDocumentFromReader(raw)
	if err != nil {
		return nil, err
	}
	return &page{doc: doc}, nil
}

func (p *page) GetTitle() string {
	return p.doc.Find("title").First().Text()
}

func (p *page) GetLinks() []string {
	var urls []string
	p.doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		url, ok := s.Attr("href")
		if ok {
			urls = append(urls, url)
		}
	})
	return urls
}
