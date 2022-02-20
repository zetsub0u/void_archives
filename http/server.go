package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zetsub0u/void_archives/archive"
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/ginrus"
	_ "github.com/zetsub0u/objcache/docs" // Need to import this to register swag instance
	ginprometheus "github.com/zsais/go-gin-prometheus"

	cors "github.com/rs/cors/wrapper/gin"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

const shutdownTimeout = 30 * time.Second

// ServerConfig TODO
type ServerConfig struct {
	Address string
	Port    int
	Version Version
}

type Version struct {
	Version string
	Commit  string
	Branch  string
}

// Server holds the information to run an HTTP Server
type Server struct {
	Config  *ServerConfig
	quit    chan struct{}
	done    chan struct{}
	engine  *gin.Engine
	version Version
	archive archive.Archive
}

// NewServer is a constructor of an HTTP Server with no logic.
func NewServer(config *ServerConfig) *Server {
	quit := make(chan struct{})
	done := make(chan struct{})
	if log.GetLevel() == log.DebugLevel {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	return &Server{
		Config:  config,
		quit:    quit,
		done:    done,
		engine:  gin.New(),
		version: config.Version,
	}
}

// WithMetrics TODO
func (s *Server) WithMetrics() *Server {
	// Monitoring
	p := ginprometheus.NewPrometheus("gin")
	// this swaps the URL with the endpoint name, to preserve cardinality
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		return c.FullPath()
	}
	p.Use(s.engine)

	return s
}

func (s *Server) WithArchive(archive archive.Archive) *Server {
	s.archive = archive
	return s
}

// Setup the server prior to start step, useful for testing.
// @title Void Archives
// @version 1.0
// @description An archive of HI3rd reference videos for MA/Abyss, trying to catalogue as much metadata as possible.
func (s *Server) Setup() {
	// Middlewares
	s.engine.Use(cors.AllowAll())

	s.engine.Use(ginrus.Ginrus(log.WithField("component", "gin"), time.RFC3339, true))
	s.engine.Use(gin.Recovery())

	// endpoints
	s.engine.GET("/version", s.GetVersion)
	s.engine.GET("/refs", s.GetRefs)

	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}

// Start the web server and serve endpoints configured in the Setup stage.
func (s *Server) Start() {
	// Server
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.Config.Address, s.Config.Port),
		Handler: s.engine,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("http: error starting http listener: %v", err)
		}
		log.Info("http: listener stopped")
	}()

	log.Infof("http: server started listening on http://%s/", srv.Addr)
	<-s.quit
	log.Printf("http: shutting down server (timeout: %s)...", shutdownTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("http: server shutdown error: ", err)
	}
	s.done <- struct{}{}
}

// Stop TODO
func (s *Server) Stop() {
	log.Info("http: stopping server...")
	s.quit <- struct{}{}
	<-s.done
	log.Info("http: server stopped.")
}
