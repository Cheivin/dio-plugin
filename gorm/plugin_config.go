package orm

import (
	core "github.com/cheivin/dio-core"
	"gorm.io/gorm"
	"strings"
)

func Gorm(options ...gorm.Option) core.PluginConfig {
	return func(d core.Dio) {
		if !d.HasProperty("gorm") {
			d.SetDefaultProperty("gorm", map[string]interface{}{
				"username": "root",
				"password": "root",
				"host":     "localhost",
				"port":     3306,
				"pool": map[string]interface{}{
					"max-idle": 0,
					"max-open": 0,
				},
				"log.level": 4,
			})
		}
		d.RegisterBean(&gormOptions{Options: options})
		d.Provide(configuration{})
	}
}

func MultiGorm(multi []string, options ...gorm.Option) core.PluginConfig {
	return func(d core.Dio) {
		if !d.HasProperty("gorm") {
			d.SetDefaultProperty("gorm", map[string]interface{}{
				"username": "root",
				"password": "root",
				"host":     "localhost",
				"port":     3306,
				"pool": map[string]interface{}{
					"max-idle": 0,
					"max-open": 0,
				},
				"log.level": 4,
			})
		}
		d.SetProperty(enableMulti, strings.Join(multi, ","))
		d.RegisterBean(&gormOptions{Options: options})
		d.Provide(configuration{})
	}
}
