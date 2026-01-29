package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chat/internal/handler"
	"chat/internal/storage"
)

func main() {
	err := storage.Connect()
	if err != nil {
		log.Panic(err)
	}

	router := http.NewServeMux()
	chatHandler := handler.NewChatHandler()

	router.HandleFunc("POST /chats/", chatHandler.AddChat)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Println("server started on port 8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panicf("listen: %s\n", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Panicf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped gracefully")
}
