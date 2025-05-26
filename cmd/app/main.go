package main

import (
	"markuzitto-bobritto/api-gateway/internal/delivery"
	"os"
	"os/signal"
	"syscall"
)

// logic here
func main() {
	server := delivery.NewServer(":8080")

	go server.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	server.Stop()
}
