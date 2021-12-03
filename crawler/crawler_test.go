package crawler

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/t0pep0/GB_best_go1/crawlerer"
	"github.com/t0pep0/GB_best_go1/crawlerer/mocks"
)

func TestNewRequester(t *testing.T) {
	r := NewRequester(10 * time.Second)
	assert.NotNil(t, r, "Error on create new requester, get nil.")
}

func TestNewCrawler(t *testing.T) {
	r := NewRequester(10 * time.Second)
	cr := NewCrawler(r)
	assert.NotNil(t, cr, "Error on create new Crawler, get nil.")
}

func TestScan(t *testing.T) {
	r := &mocks.Requester{}
	pg := &mocks.Page{}
	cr := NewCrawler(r)

	ctx := context.Background()
	var depth int64 = 3
	var start int64 = 1
	url := "https://testStartURL.com"

	exVisitedURLCount := 4
	exTitle := "Test title"
	exLinks := []string{"https://childURL1.com", "https://childURL2.com", "https://childURL3.com"}

	pg.On("GetTitle").Return(exTitle)
	pg.On("GetLinks").Return(exLinks)

	r.On("Get", ctx, url).Return(pg, nil)
	r.On("Get", ctx, exLinks[0]).Return(pg, nil)
	r.On("Get", ctx, exLinks[1]).Return(pg, nil)
	r.On("Get", ctx, exLinks[2]).Return(pg, nil)

	go cr.Scan(ctx, url, url, &depth, start)

	var maxResult, maxErrors = 10, 5
	doFor := true
	for doFor {
		// gosimple
		// delete select statement
		// S1000: should use a simple channel send/receive instead of `select` with a single case (gosimple)
		msg := <-cr.ChanResult()
		if msg.Err != nil {
			maxErrors--
			if maxErrors <= 0 {
				doFor = false
			}
		} else if len(msg.Info) > 0 {
			doFor = false
		} else {
			maxResult--
			if maxResult <= 0 {
				doFor = false
			}
		}
	}

	assert.Equal(t, exVisitedURLCount, len(cr.visited), "Not equal.\n Expected: %v \n Got: %v \n", exVisitedURLCount, len(cr.visited))
}

func TestToChanResult(t *testing.T) {
	r := NewRequester(1)
	cr := NewCrawler(r)

	crawlRes := crawlerer.CrawlResult{
		Err:   nil,
		Info:  "",
		Title: "Test title",
		Url:   "https://testCurrentURL.com",
	}

	go cr.ToChanResult(crawlRes)

	gotRes := <-cr.res
	assert.Equal(t, crawlRes, gotRes, "Not equal.\n Expected: %v \n Got: %v \n", crawlRes, gotRes)
}

func TestChanResult(t *testing.T) {
	r := NewRequester(10 * time.Second)
	cr := NewCrawler(r)

	crawlRes := crawlerer.CrawlResult{
		Err:   nil,
		Info:  "",
		Title: "Test title",
		Url:   "https://testCurrentURL.com",
	}

	go func() {
		cr.res <- crawlRes
	}()

	gotRes := <-cr.ChanResult()
	assert.Equal(t, crawlRes, gotRes, "Not equal.\n Expected: %v \n Got: %v \n", crawlRes, gotRes)
}

func TestGet(t *testing.T) {
	http.HandleFunc("/home/", func(w http.ResponseWriter, r *http.Request) {
		type Link struct {
			URL   string
			Title string
		}

		type templateData struct {
			Links []Link
		}

		links := []Link{
			{URL: "https://childURL1.com", Title: "link 1"},
			{URL: "https://childURL2.com", Title: "link 2"},
			{URL: "https://childURL3.com", Title: "link 3"},
		}

		tmpl, err := template.ParseFiles("../crawlerer/mocks/home.page.tmpl")
		if err != nil {
			t.Fatalf("Error on parse template: %v \n", err)
		}
		err = tmpl.Execute(w, &templateData{Links: links})
		if err != nil {
			t.Fatalf("Error on execute template: %v \n", err)
		}
	})

	addr := "localhost:8080"
	url := "http://" + addr + "/home/"
	server := &http.Server{Addr: addr, Handler: nil}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			t.Errorf("Error on start server: %v \n", err)
			return
		}
	}()

	time.Sleep(3 * time.Second)

	r := NewRequester(3 * time.Second)
	ctx := context.Background()

	pg, err := r.Get(ctx, url)
	if err != nil {
		t.Fatalf("Error on get page: %v \n", err)
	}

	assert.NotNil(t, pg, "Get nil page.")
}

func TestNewPage(t *testing.T) {
	type Link struct {
		URL   string
		Title string
	}

	type templateData struct {
		Links []Link
	}

	links := []Link{
		{URL: "https://childURL1.com", Title: "link 1"},
		{URL: "https://childURL2.com", Title: "link 2"},
		{URL: "https://childURL3.com", Title: "link 3"},
	}

	tmpl, err := template.ParseFiles("../crawlerer/mocks/home.page.tmpl")
	if err != nil {
		t.Fatalf("Error on parse template: %v \n", err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, &templateData{Links: links})
	// ineffassign
	// add error handling
	// ineffectual assignment to err (ineffassign)
	if err != nil {
		t.Fatalf("Error on execute template: %v \n", err)
	}

	pg, err := NewPage(&buffer)
	// ineffassign
	// add error handling
	// ineffectual assignment to err (ineffassign)
	if err != nil {
		t.Fatalf("Error on create new page struct: %v \n", err)
	}

	assert.NotNil(t, pg, "Get nil page.")
}

func TestGetTitle(t *testing.T) {
	type Link struct {
		URL   string
		Title string
	}

	type templateData struct {
		Links []Link
	}

	links := []Link{
		{URL: "https://childURL1.com", Title: "link 1"},
		{URL: "https://childURL2.com", Title: "link 2"},
		{URL: "https://childURL3.com", Title: "link 3"},
	}

	tmpl, err := template.ParseFiles("../crawlerer/mocks/home.page.tmpl")
	if err != nil {
		t.Fatalf("Error on parse template: %v \n", err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, &templateData{Links: links})
	// ineffassign
	// add error handling
	// ineffectual assignment to err (ineffassign)
	if err != nil {
		t.Fatalf("Error on execute template: %v \n", err)
	}

	pg, err := NewPage(&buffer)
	// ineffassign
	// add error handling
	// ineffectual assignment to err (ineffassign)
	if err != nil {
		t.Fatalf("Error on create new page struct: %v \n", err)
	}

	exTitle := "Home page"
	gotTitle := pg.GetTitle()

	assert.Equal(t, exTitle, gotTitle, "Not equal.\n Expected: %v \n Got: %v \n", exTitle, gotTitle)
}

func TestGetLinks(t *testing.T) {
	type Link struct {
		URL   string
		Title string
	}

	type templateData struct {
		Links []Link
	}

	links := []Link{
		{URL: "https://childURL1.com", Title: "link 1"},
		{URL: "https://childURL2.com", Title: "link 2"},
		{URL: "https://childURL3.com", Title: "link 3"},
	}

	tmpl, err := template.ParseFiles("../crawlerer/mocks/home.page.tmpl")
	if err != nil {
		t.Fatalf("Error on parse template: %v \n", err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, &templateData{Links: links})
	// ineffassign
	// add error handling
	// ineffectual assignment to err (ineffassign)
	if err != nil {
		t.Fatalf("Error on execute template: %v \n", err)
	}

	pg, err := NewPage(&buffer)
	// ineffassign
	// add error handling
	// ineffectual assignment to err (ineffassign)
	if err != nil {
		t.Fatalf("Error on create new page struct: %v \n", err)
	}

	exLinks := []string{"https://childURL1.com", "https://childURL2.com", "https://childURL3.com"}
	gotLinks := pg.GetLinks()

	assert.Equal(t, exLinks, gotLinks, "Not equal.\n Expected: %v \n Got: %v \n", exLinks, gotLinks)
}
