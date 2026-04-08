package router

import (
	"embed"
	"net/http"
	"path"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/controller"
	"github.com/QuantumNous/new-api/middleware"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func shouldServeSPA(c *gin.Context) bool {
	if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
		return false
	}

	requestPath := c.Request.URL.Path
	if requestPath == "" {
		requestPath = "/"
	}

	if strings.HasPrefix(requestPath, "/v1") ||
		strings.HasPrefix(requestPath, "/api") ||
		strings.HasPrefix(requestPath, "/assets") {
		return false
	}

	if path.Ext(requestPath) != "" {
		return false
	}

	segments := strings.Split(strings.Trim(requestPath, "/"), "/")
	for _, segment := range segments {
		segment = strings.ToLower(strings.TrimSpace(segment))
		if segment == "" {
			continue
		}
		if strings.HasPrefix(segment, ".") {
			return false
		}
		if segment == "vendor" || segment == "phpunit" || segment == "containers" {
			return false
		}
	}

	accept := strings.ToLower(c.GetHeader("Accept"))
	return accept == "" ||
		strings.Contains(accept, "text/html") ||
		strings.Contains(accept, "application/xhtml+xml") ||
		strings.Contains(accept, "*/*")
}

func SetWebRouter(router *gin.Engine, buildFS embed.FS, indexPage []byte) {
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.GlobalWebRateLimit())
	router.Use(middleware.Cache())
	router.Use(static.Serve("/", common.EmbedFolder(buildFS, "web/dist")))
	router.NoRoute(func(c *gin.Context) {
		c.Set(middleware.RouteTagKey, "web")
		if !shouldServeSPA(c) {
			controller.RelayNotFound(c)
			return
		}
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexPage)
	})
}
