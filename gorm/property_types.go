package orm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/url"
	"time"
)

type gormOptions struct {
	Options []gorm.Option
}

func (c *gormOptions) BeanName() string {
	return "gormOptions"
}

// DBProperty 数据库配置
type DBProperty struct {
	Username   string `value:"username"`
	Password   string `value:"password"`
	Host       string `value:"host"`
	Port       int    `value:"port"`
	Database   string `value:"database"`
	Parameters string `value:"parameters"`
}

func (p *DBProperty) parseParameters() {
	if p.Parameters == "" {
		return
	}
	_, err := url.ParseQuery(p.Parameters)
	if err != nil {
		panic(err)
	}
}

func (p *DBProperty) Merge(property DBProperty) {
	if p.Username == "" && property.Username != "" {
		p.Username = property.Username
	}
	if p.Password == "" && property.Password != "" {
		p.Password = property.Password
	}
	if p.Host == "" && property.Host != "" {
		p.Host = property.Host
	}
	if p.Port == 0 && property.Port > 0 {
		p.Port = property.Port
	}
	if p.Database == "" && property.Database != "" {
		p.Database = property.Database
	}
	if p.Parameters == "" && property.Parameters != "" {
		p.Parameters = property.Parameters
	}
}

// PoolProperty 连接池配置
type PoolProperty struct {
	MaxIdle     int           `value:"pool.max-idle"`
	MaxOpen     int           `value:"pool.max-open"`
	MaxLifeTime time.Duration `value:"pool.max-life-time"`
	MaxIdleTime time.Duration `value:"pool.max-idle-time"`
}

func (p *PoolProperty) Merge(property PoolProperty) {
	if p.MaxIdle == 0 && property.MaxIdle > 0 {
		p.MaxIdle = property.MaxIdle
	}
	if p.MaxOpen == 0 && property.MaxOpen > 0 {
		p.MaxOpen = property.MaxOpen
	}
	if p.MaxLifeTime == 0 && property.MaxLifeTime > 0 {
		p.MaxLifeTime = property.MaxLifeTime
	}
	if p.MaxIdleTime == 0 && property.MaxIdleTime > 0 {
		p.MaxIdleTime = property.MaxIdleTime
	}
}

// LogProperty 日志配置
type LogProperty struct {
	Level                     int           `value:"log.level"`
	SlowThreshold             time.Duration `value:"log.slow-log"`
	IgnoreRecordNotFoundError bool          `value:"log.ignore-notfound"`
}

func (p LogProperty) LogLevel() logger.LogLevel {
	if p.Level <= 0 || p.Level > 4 {
		return logger.Info
	}
	return logger.LogLevel(p.Level)
}

func (p LogProperty) Merge(property LogProperty) {
	if p.Level == 0 && property.Level > 0 {
		p.Level = property.Level
	}
	if p.SlowThreshold == 0 && property.SlowThreshold > 0 {
		p.SlowThreshold = property.SlowThreshold
	}
}
