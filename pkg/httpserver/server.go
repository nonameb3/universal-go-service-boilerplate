package httpserver

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	App  *fiber.App
	port int
}

func New(port int) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	})

	return &Server{
		App:  app,
		port: port,
	}
}

func (s *Server) Start() error {
	return s.App.Listen(":" + strconv.Itoa(s.port))
}

func (s *Server) StartWithGracefulShutdown() {
	go func() {
		if err := s.App.Listen(":" + strconv.Itoa(s.port)); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.App.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}
