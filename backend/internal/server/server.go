package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crimsonn/zm_server/internal/database"
	"github.com/crimsonn/zm_server/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Host string
	Port int
}

type Server struct {
	Config   Config
	Router   *gin.Engine
	Database *sqlx.DB
	Logger   *logrus.Entry
	handlers *ServerHandlers
}

func NewServer(config Config) *Server {
	db, err := database.NewPostgres()
	if err != nil {
		log.Fatalf("failed to initialize postgres: %v", err)
	}

	server := &Server{
		Config:   config,
		Router:   gin.Default(),
		Database: db,
		Logger:   logger.NewComponentLogger("api-server"),
		handlers: NewServerHandlers(db),
	}
	server.setupMiddlewares()
	server.handlers.RegisterRoutes(server.Router)
	return server
}

func (s *Server) setupMiddlewares() {
	s.Router.Use(gin.Recovery())
	s.Router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - %s \"%s %s %s\" %d %s %s\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.ErrorMessage,
			)
		},
		Output: s.Logger.Writer(),
	}))

}

func (s *Server) Start(done chan bool) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		addr := fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port)
		s.Logger.Infof("starting server on %s", addr)
		if err := s.Router.Run(addr); err != nil {
			s.Logger.Errorf("failed to start server: %v", err)
		}
	}()

	<-quit
	s.Logger.Info("shutdown signal received...")
	close(done)
	time.Sleep(1 * time.Second)
	s.Logger.Info("server shutting down...")
	if err := s.Database.Close(); err != nil {
		s.Logger.Errorf("failed to close database connection: %v", err)
		return err
	}
	s.Logger.Info("database connection closed")
	s.Logger.Info("server shutdown complete")
	return nil
}
