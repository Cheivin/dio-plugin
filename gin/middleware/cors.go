package middleware

import (
	"github.com/cheivin/dio-core/system"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// WebCors 跨域
type WebCors struct {
	Web              *gin.Engine   `aware:"web"`
	Log              *system.Log   `aware:""`
	Origins          string        `value:"app.web.cors.origin"`
	Methods          string        `value:"app.web.cors.method"`
	Headers          string        `value:"app.web.cors.header"`
	AllowCredentials bool          `value:"app.web.cors.allow-credentials"`
	ExposeHeaders    string        `value:"app.web.cors.expose-header"`
	MaxAge           time.Duration `value:"app.web.cors.max-age"` // 过期时间,单位秒
	config           cors.Config
}

func (w *WebCors) BeanConstruct() {
	w.config = cors.DefaultConfig()
	if w.Origins != "" {
		w.config.AllowOrigins = strings.Split(w.Origins, ",")
	} else {
		w.config.AllowAllOrigins = true
	}
	if w.Methods != "" {
		w.config.AllowMethods = strings.Split(w.Methods, ",")
	}
	if w.Headers != "" {
		w.config.AllowHeaders = strings.Split(w.Headers, ",")
	}
	w.config.AllowCredentials = w.AllowCredentials
	if w.ExposeHeaders != "" {
		w.config.ExposeHeaders = strings.Split(w.ExposeHeaders, ",")
	}

	if w.MaxAge.Seconds() > 0 {
		w.config.MaxAge = w.MaxAge
	}
}

func (w *WebCors) AfterPropertiesSet() {
	w.Web.Use(cors.New(w.config))
}
