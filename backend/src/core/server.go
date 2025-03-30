package core

import (
	"context"
	"net/http"
	"slices"
	"teirxserver/src/txlog"

	"github.com/gin-gonic/gin"
)

type GinLogForwarder struct{}

func (g GinLogForwarder) Write(p []byte) (n int, err error) {
	txlog.TxLogInfo(string(p))
	return n, nil
}

type Server struct {
	httpServer *http.Server
}

func NewHTTPServer() *Server {
	gin.DisableConsoleColor()
	// gin.DefaultWriter = GinLogForwarder{}
	router := gin.Default()
	router.Use(corsMiddleware())
	router.Static("/images", "./images")
	RegisterRoutes(router)
	return &Server{
		httpServer: &http.Server{
			Addr:    "localhost:8080",
			Handler: router,
		},
	}
}

func (s *Server) Start() {
	go func() {
		txlog.TxLogInfo("Starting HTTP server: %s", s.httpServer.Addr)
		err := s.httpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			txlog.TxLogError("Http server closed: %s", err)
		}
	}()
}

func (s *Server) Stop() {
	err := s.httpServer.Shutdown(context.Background())
	if err != nil {
		txlog.TxLogError("Error shutting down HTTP server: %s", err.Error())
	}
}

// CORS middleware function definition
func corsMiddleware() gin.HandlerFunc {

	var allowedOrigins = []string{
		"http://localhost:3000",
	}

	// Return the actual middleware handler function
	return func(c *gin.Context) {

		origin := c.Request.Header.Get("Origin")

		if slices.Contains(allowedOrigins, origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		}

		// Handle preflight OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent) 
			return
		}

		// Continue with the request
		c.Next()
	}
}
