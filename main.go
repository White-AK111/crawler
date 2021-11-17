package main

import (
	"context"
	"github.com/t0pep0/GB_best_go1/config"
	"github.com/t0pep0/GB_best_go1/crawler"
	"github.com/t0pep0/GB_best_go1/crawlerer"
	"log"
	"os"
	"os/signal"
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
	ctxCurr, cancelCurr := context.WithCancel(ctxTime)

	r := crawler.NewRequester(cfg.App.TimeoutRequest * time.Second)
	cr := crawler.NewCrawler(r)

	go cr.Scan(ctxCurr, cfg.App.URL, cfg.App.URL, &cfg.App.MaxDepth, 1) // Запускаем краулер в отдельной рутине
	go processResult(ctxCurr, cancelCurr, cr, cfg)                      // Обрабатываем результаты в отдельной рутине

	sigCh := make(chan os.Signal)                                          // Создаем канал для приема сигналов
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2) // Подписываемся на сигнал SIGINT
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
				cfg.ChangeMaxDepth(cfg.App.DeltaDepth)
				log.Printf("Depth increment set to: %d\n", cfg.App.DeltaDepth)
			case syscall.SIGUSR2:
				log.Println("Catch SIGUSR2")
				// Decrement depth while catch SIGUSR2
				cfg.ChangeMaxDepth(-cfg.App.DeltaDepth)
				log.Printf("Depth decrement set to: -%d\n", cfg.App.DeltaDepth)
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
				log.Printf("crawler result return err: %s\n", msg.Err.Error())
				if maxErrors <= 0 {
					cancel()
					log.Println("Exit by MaxErrors")
					return
				}
			} else if len(msg.Info) > 0 {
				cancel()
				log.Printf("Exit by: %s\n", msg.Info)
				return
			} else {
				maxResult--
				log.Printf("crawler result: [url: %s] Title: %s\n", msg.Url, msg.Title)
				if maxResult <= 0 {
					cancel()
					log.Println("Exit by MaxResults")
					return
				}
			}
		}
	}
}
