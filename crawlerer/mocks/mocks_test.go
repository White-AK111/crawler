package mocks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/t0pep0/GB_best_go1/crawlerer"
)

// TestCrawler contracts of Crawler interface
func TestCrawler(t *testing.T) {
	cr := &Crawler{}

	ctx := context.Background()
	crawlRes := crawlerer.CrawlResult{
		Err:   nil,
		Info:  "",
		Title: "Test title",
		Url:   "https://testCurrentURL.com",
	}

	var depth int64 = 3
	var start int64 = 1
	url := "https://testStartURL.com"

	// ineffassign
	// delete resChan ineffectual assignment
	// ineffectual assignment to resChan (ineffassign)

	cr.On("Scan", ctx, url, url, &depth, start)
	cr.On("ChanResult").Return(make(<-chan crawlerer.CrawlResult))
	cr.On("ToChanResult", crawlRes)

	cr.Scan(ctx, url, url, &depth, start)

	// add assignment ":="
	resChan := cr.ChanResult()
	cr.ToChanResult(crawlRes)
	assert.NotNil(t, resChan, "Got nil result chan.")
	require.NotZero(t, cr, "Error on check contract in crawler interface.")
}

// TestRequester contracts of Requester interface
func TestRequester(t *testing.T) {
	r := &Requester{}
	pg := &Page{}

	ctx := context.Background()
	url := "https://testCurrentURL.com"

	r.On("Get", ctx, url).Return(pg, nil)

	res, err := r.Get(ctx, url)
	require.NoError(t, err, "Error on check contract in requester interface.")
	assert.Equal(t, pg, res, "Not equal.\n Expected: %v \n Got: %v \n", pg, res)
}

// TestPage contracts of Page interface
func TestPage(t *testing.T) {
	pg := &Page{}

	exTitle := "Test title"
	exLinks := []string{"https://childURL1.com", "https://childURL2.com", "https://childURL3.com"}

	pg.On("GetTitle").Return(exTitle)
	pg.On("GetLinks").Return(exLinks)

	title := pg.GetTitle()
	links := pg.GetLinks()
	assert.Equal(t, exTitle, title, "Not equal.\n Expected: %v \n Got: %v \n", exTitle, title)
	assert.Equal(t, exLinks, links, "Not equal.\n Expected: %v \n Got: %v \n", exLinks, links)
}
