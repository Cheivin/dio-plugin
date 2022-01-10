package middleware

import (
	"github.com/cheivin/di"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type corsConfig struct {
	Origins          string        `value:"origin"`
	Methods          string        `value:"method"`
	Headers          string        `value:"header"`
	AllowCredentials bool          `value:"allow-credentials"`
	ExposeHeaders    string        `value:"expose-header"`
	MaxAge           time.Duration `value:"max-age"` // 过期时间,单位秒
}

// WebCors 跨域
type WebCors struct {
	Web *gin.Engine `aware:"web"`
}

func (w *WebCors) AfterPropertiesSet(container di.DI) {
	cfg := container.LoadProperties("app.web.cors.", corsConfig{}).(corsConfig)
	corsCfg := cors.DefaultConfig()
	if cfg.Origins != "" {
		corsCfg.AllowOrigins = strings.Split(cfg.Origins, ",")
	} else {
		corsCfg.AllowAllOrigins = true
	}
	if cfg.Methods != "" {
		corsCfg.AllowMethods = strings.Split(cfg.Methods, ",")
	}
	if cfg.Headers != "" {
		corsCfg.AllowHeaders = strings.Split(cfg.Headers, ",")
	}
	corsCfg.AllowCredentials = cfg.AllowCredentials
	if cfg.ExposeHeaders != "" {
		corsCfg.ExposeHeaders = strings.Split(cfg.ExposeHeaders, ",")
	}
	if cfg.MaxAge.Seconds() > 0 {
		corsCfg.MaxAge = cfg.MaxAge
	}

	w.Web.Use(cors.New(corsCfg))
}
