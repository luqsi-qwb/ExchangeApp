package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/kuqsi/exchangeapp/config"
	"github.com/kuqsi/exchangeapp/route"
)

func main() {
	config.InitConfig()

	r := route.SetupRouter()

	port := config.AppConfig.App.Port

	if port == "" {
		port = ":8080"
	}

	//实现一下优雅退出
	src := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		err := src.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen err,err is %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown,Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := src.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown %v", err)
	}
	log.Println("server exiting")

}
