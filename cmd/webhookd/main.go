package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LYH2263/go-webhookhub"
	"github.com/LYH2263/go-webhookhub/internal/api"
)

func main() {
	addr := flag.String("addr", ":8090", "HTTP 监听地址")
	web := flag.String("web", "web", "静态管理页目录")
	persist := flag.String("persist", "", "端点快照 JSON 路径（可选）")
	journal := flag.String("journal", "", "投递日志追加文件（可选）")
	timeout := flag.Duration("http-timeout", 10*time.Second, "默认投递超时")
	maxTry := flag.Int("max-attempts", 3, "默认最大尝试次数")
	flag.Parse()

	opts := []webhookhub.Option{
		webhookhub.WithHTTPTimeout(*timeout),
		webhookhub.WithMaxAttempts(*maxTry),
		webhookhub.WithAllowHTTP(true),
	}
	if *persist != "" {
		opts = append(opts, webhookhub.WithPersistPath(*persist))
	}
	if *journal != "" {
		opts = append(opts, webhookhub.WithJournalPath(*journal))
	}

	h := webhookhub.New(opts...)
	defer h.Close()

	srv := api.New(h, api.Options{WebDir: *web, AllowCORS: true})
	httpSrv := &http.Server{
		Addr:              *addr,
		Handler:           srv,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("webhookd 监听 %s", *addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	_ = httpSrv.Close()
}
