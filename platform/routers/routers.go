package routers

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/adiubaidah/rfid-syafiiyah/pkg/config"
	"github.com/gin-gonic/gin"
)

// Router contains the functions of http handler to clean payloads and pass it the service
type Router interface {
	Serve()
}

// Route data will be registered to http listener
type Route struct {
	Method      string
	Path        string
	Handle      gin.HandlerFunc
	MiddleWares gin.HandlersChain
}

type routing struct {
	address string
	routers []Route
}

// NewRouting is for creating new routing
func NewRouting(address string, routers []Route) Router {
	return &routing{
		address,
		routers,
	}
}

func (r *routing) Serve() {
	ginRouter := gin.New()
	ginRouter.Use(gin.Logger())
	ginRouter.Use(gin.Recovery())
	ginRouter.Use(CORSHandler)
	ginRouter.Static("/photo", config.PathPhoto)
	ginRouter.Handle(http.MethodGet, "/ping", HealthCheck)

	for _, router := range r.routers {
		if router.MiddleWares == nil {
			ginRouter.Handle(router.Method, router.Path, router.Handle)
		} else {
			var handlers []gin.HandlerFunc
			for _, middle := range router.MiddleWares {
				handlers = append(handlers, middle)
			}
			handlers = append(handlers, router.Handle)

			ginRouter.Handle(router.Method, router.Path, handlers...)
		}
	}

	err := ginRouter.Run(r.address)
	if err != nil {
		panic(err)
	}
	fmt.Printf("started at %s", r.address)
	ginRouter.Run(r.address)
}

// CORSHandler handles requests with unsupported HTTP methods.
func CORSHandler(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	if origin != "" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	}
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}

// HealthCheck handles the HTTP request for health checking the service.
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"description": runtime.NumGoroutine(),
	})
}
