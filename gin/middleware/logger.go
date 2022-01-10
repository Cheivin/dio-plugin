package middleware

import (
	"github.com/cheivin/di"
	"github.com/cheivin/dio-core/system"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

type logConfig struct {
	Skips     string `value:"skip-path"`
	TraceName string `value:"trace-name"`
}

func (c logConfig) SkipPaths() map[string]struct{} {
	skipPaths := strings.Split(c.Skips, ",")
	skipMap := make(map[string]struct{}, len(skipPaths))
	for _, path := range skipPaths {
		if path != "" {
			skipMap[path] = struct{}{}
		}
	}
	return skipMap
}

type (
	WebTracert interface {
		Trace(c *gin.Context, traceName string) string
	}

	// WebLogger 日志
	WebLogger struct {
		Log     *system.Log `aware:""`
		Web     *gin.Engine `aware:"web"`
		Tracert WebTracert  `aware:"omitempty"`
		config  logConfig
		skip    map[string]struct{}
	}

	defaultWebTracert struct {
		UUID uuid.UUID
	}
)

func newDefaultWebTracert() WebTracert {
	return &defaultWebTracert{
		UUID: uuid.NewV4(),
	}
}

func (t defaultWebTracert) Trace(c *gin.Context, traceName string) string {
	reqId := c.GetHeader(traceName)
	if reqId == "" {
		reqId = t.UUID.String()
		c.Header(traceName, reqId)
	}
	c.Set(traceName, reqId)
	return reqId
}

func (w *WebLogger) AfterPropertiesSet(container di.DI) {
	w.Log = w.Log.WithOptions(zap.WithCaller(false))
	if w.Tracert == nil {
		w.Tracert = newDefaultWebTracert()
	}

	w.config = container.LoadProperties("app.web.log.", logConfig{}).(logConfig)
	w.skip = w.config.SkipPaths()
	w.Web.Use(w.log)
}

func (w *WebLogger) log(c *gin.Context) {
	defer func() {
		// 此处recover用于处理顶层log中间件写出日志panic
		if r := recover(); r != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()
	// 开始时间
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery

	// 跟踪id
	w.Tracert.Trace(c, w.config.TraceName)

	// 处理请求
	c.Next()

	// 判断是否过滤路径
	for skipPath := range w.skip {
		if strings.HasPrefix(path, skipPath) {
			return
		}
	}

	// 记录日志
	timeStamp := time.Now()
	if raw != "" {
		path = path + "?" + raw
	}
	params := []interface{}{
		"TimeStamp", timeStamp,
		"Cost", timeStamp.Sub(start).String(),
		"ClientIP", c.ClientIP(),
		"Method", c.Request.Method,
		"StatusCode", c.Writer.Status(),
		"Path", path,
		"BodySize", c.Writer.Size(),
	}
	errMsg := c.Errors.Last()
	if errMsg != nil {
		params = append(params, "ErrorMessage", errMsg)
		w.Log.Error(c, "gin-http", params...)
	} else {
		w.Log.Info(c, "gin-http", params...)
	}
}
