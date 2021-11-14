package main

import (
	"context"
	"github.com/t0pep0/GB_best_go1/config"
	"github.com/t0pep0/GB_best_go1/crawler"
	"github.com/t0pep0/GB_best_go1/crawlerer"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {
	// Load config file
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("Can't load configuration file: %s", err)
	}

	ctxTime, cancelTimeout := context.WithTimeout(context.Background(), cfg.App.TimeoutApp*time.Second)
	ctxCurr, cancelCurr := crawlerStart(ctxTime, cfg)

	sigCh := make(chan os.Signal)                         // Создаем канал для приема сигналов
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGUSR1) // Подписываемся на сигнал SIGINT
	for {
		select {
		case <-ctxCurr.Done(): // Если всё завершили - выходим
			return
		case sig := <-sigCh:
			switch sig {
			case syscall.SIGINT:
				// Close while catch SIGINT
				log.Println("Catch SIGINT")
				cancelCurr()
				return
			case syscall.SIGUSR1:
				log.Println("Catch SIGUSR1")
				// Increment depth while catch SIGUSR1
				atomic.AddUint64(&cfg.App.MaxDepth, cfg.App.DeltaDepth)
				log.Printf("Depth increment set to: %d\n", cfg.App.DeltaDepth)
			default:
				log.Println("Catch other signal")
			}
		case <-ctxTime.Done():
			// Exit by application timeout
			{
				log.Println("Exit by timeout")
				cancelCurr()
				cancelTimeout()
			}
		}
	}
}

// processResult get results of process
func processResult(ctx context.Context, cancel context.CancelFunc, cr crawlerer.Crawler, cfg *config.Config) {
	var maxResult, maxErrors = cfg.App.MaxResults, cfg.App.MaxErrors
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-cr.ChanResult():
			if msg.Err != nil {
				maxErrors--
				log.Printf("crawler result return err: %s Depth: %d\n", msg.Err.Error(), cfg.App.MaxDepth)
				if maxErrors <= 0 {
					cancel()
					log.Println("Exit by MaxErrors")
					return
				}
			} else {
				maxResult--
				log.Printf("crawler result: [url: %s] Title: %s Depth: %d\n", msg.Url, msg.Title, cfg.App.MaxDepth)
				if maxResult <= 0 {
					cancel()
					log.Println("Exit by MaxResults")
					return
				}
			}
		}
	}
}

// crawlerStart starter for crawler
func crawlerStart(ctxMain context.Context, cfg *config.Config) (ctxCurr context.Context, cancelCurr context.CancelFunc) {
	ctxCurr, cancelCurr = context.WithCancel(ctxMain)
	r := crawler.NewRequester(cfg.App.TimeoutRequest * time.Second)
	cr := crawler.NewCrawler(r)

	go cr.Scan(ctxCurr, cfg.App.URL, cfg.App.URL, uint(cfg.App.MaxDepth)) // Запускаем краулер в отдельной рутине
	go processResult(ctxCurr, cancelCurr, cr, cfg)                        // Обрабатываем результаты в отдельной рутине

	return ctxCurr, cancelCurr
}
