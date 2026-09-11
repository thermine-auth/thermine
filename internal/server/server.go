package server

import (
	"github.com/gin-gonic/gin"

	"xermess/internal/config"
)

func New(cfg config.Config) *gin.Engine {
	r := gin.New()

	r.Use(
		gin.Logger(),
		gin.Recovery(),
		RequestID(),
		SecurityHeaders(),
		CORS(cfg.CORSOrigins),
	)

	// Client IPs come from the connection only. Set the real proxy CIDRs here
	// once this runs behind a load balancer, or ClientIP() can be spoofed
	// through X-Forwarded-For.
	_ = r.SetTrustedProxies(nil)

	registerRoutes(r)

	return r
}
