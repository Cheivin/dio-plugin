package middleware

import (
	"github.com/cheivin/dio-core/system"
	"github.com/cheivin/dio-plugin/gin/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	WebErrorHandler interface {
		OnError(c *gin.Context, err errors.Error)
	}

	// WebRecover 全局异常
	WebRecover struct {
		Web          *gin.Engine     `aware:"web"`
		Log          *system.Log     `aware:""`
		ErrorHandler WebErrorHandler `aware:"omitempty"`
	}

	defaultWebErrorHandler struct {
	}
)

func newDefaultWebErrorHandler() WebErrorHandler {
	return &defaultWebErrorHandler{}
}

func (d defaultWebErrorHandler) OnError(c *gin.Context, err errors.Error) {
	c.String(err.Code, err.Error())
}

func (w *WebRecover) AfterPropertiesSet() {
	if w.ErrorHandler == nil {
		w.ErrorHandler = newDefaultWebErrorHandler()
	}

	w.Web.NoRoute(w.noRoute)
	w.Web.NoMethod(w.noMethod)
	w.Web.Use(w.recover)
}

func (w *WebRecover) noRoute(c *gin.Context) {
	_ = c.Error(errors.NoRoute.Cause())
	w.ErrorHandler.OnError(c, errors.NoRoute)
}

func (w *WebRecover) noMethod(c *gin.Context) {
	_ = c.Error(errors.NoMethod.Cause())
	w.ErrorHandler.OnError(c, errors.NoMethod)
}

func (w *WebRecover) recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case errors.Error:
				err := r.(errors.Error)
				_ = c.Error(&err)
				w.ErrorHandler.OnError(c, err)
			case error:
				e := r.(error)
				err := errors.Err(http.StatusInternalServerError, e.Error(), e)
				_ = c.Error(&err)
				w.ErrorHandler.OnError(c, err)
			case string:
				err := errors.ErrMsg(http.StatusInternalServerError, r.(string))
				_ = c.Error(&err)
				w.ErrorHandler.OnError(c, err)
			default:
				w.Log.Error(c, "Web server panic", "panic", r)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.Abort()
		}
	}()
	//加载完 defer recover，继续后续接口调用
	c.Next()
}
