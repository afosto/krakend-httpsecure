package gin

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/config"
	"github.com/unrolled/secure"

	"github.com/devopsfaith/krakend-httpsecure"
)

var errNoConfig = errors.New("no config present for the httpsecure module")

// Register registers the secure middleware into the gin engine
func Register(cfg config.ExtraConfig, engine *gin.Engine) error {
	opt, ok := httpsecure.ConfigGetter(cfg).(secure.Options)
	if !ok {
		return errNoConfig
	}
	engine.Use(secureMw(opt))
	return nil
}

// NewSecureMw creates a secured middleware for the gin engine
func NewSecureMw(cfg config.ExtraConfig) gin.HandlerFunc {
	opt, ok := httpsecure.ConfigGetter(cfg).(secure.Options)
	if !ok {
		return func(c *gin.Context) {}
	}

	return secureMw(opt)
}

// secureMw creates a secured middleware for the gin engine
func secureMw(opt secure.Options) gin.HandlerFunc {
	secureMiddleware := secure.New(opt)

	return func(c *gin.Context) {

		if isHealthEndpoint(c.Request.URL.Path) {
			return
		}

		err := secureMiddleware.Process(c.Writer, c.Request)

		if err != nil {
			c.Abort()
			return
		}

		if status := c.Writer.Status(); status > 300 && status < 399 {
			c.Abort()
		}
	}
}

func isHealthEndpoint(path string) bool {
	if path == "/healthz" || path == "/_ah/health" || path == "/__health" {
		return true
	}
	return false
}
