package web

import (
	core "github.com/cheivin/dio-core"
	"github.com/cheivin/dio-plugin/gin/middleware"
)

const defaultTraceName = core.DefaultTraceName

func GinWeb(useLogger bool, options ...ginOption) core.PluginConfig {
	return func(d core.Dio) {
		if !d.HasProperty("app.port") {
			d.SetDefaultPropertyMap(map[string]interface{}{
				"app.port": 8080,
			})
		}
		d.Provide(ginContainer{})
		if useLogger {
			if !d.HasProperty("app.web.log") {
				d.SetDefaultProperty("app.web.log", map[string]interface{}{
					"skip-path":  "",
					"trace-name": defaultTraceName,
				})
			}
			d.Provide(middleware.WebLogger{})
		}
		d.Provide(middleware.WebRecover{})
		for _, option := range options {
			option(d)
		}
		d.Provide(Controller{})
	}
}

type ginOption func(core.Dio)

func WithCors(useCors bool) ginOption {
	return func(d core.Dio) {
		if useCors {
			if !d.HasProperty("app.web.cors") {
				d.SetDefaultProperty("app.web.cors", map[string]interface{}{
					"origin":            "",
					"method":            "",
					"header":            "",
					"allow-credentials": true,
					"expose-header":     "",
					"max-age":           "12h",
				})
			}
			d.Provide(middleware.WebCors{})
		}
	}
}

func WithTracert(tracert middleware.WebTracert) ginOption {
	return func(d core.Dio) {
		d.Provide(tracert)
	}
}

func WithErrorHandler(errorHandler middleware.WebErrorHandler) ginOption {
	return func(d core.Dio) {
		d.Provide(errorHandler)
	}
}
