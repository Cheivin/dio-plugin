package orm

import (
	"fmt"
	"github.com/cheivin/di"
	"github.com/cheivin/dio-core/system"
	"github.com/cheivin/dio-plugin/gorm/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

type configuration struct {
	log                 *system.Log
	opts                []gorm.Option
	defaultDBProperty   DBProperty
	defaultPoolProperty PoolProperty
	defaultLogProperty  LogProperty
}

func (c *configuration) BeanName() string {
	return "gormConfiguration"
}

func (c *configuration) BeanConstruct(container di.DI) {
	defaultPrefix := "gorm."
	// 系统日志
	bean, _ := container.GetByType(system.Log{})
	c.log = bean.(*system.Log)
	// 配置信息
	if options, ok := container.GetByType(gormOptions{}); ok {
		c.opts = options.(*gormOptions).Options
	}
	// db连接信息
	c.defaultDBProperty = container.LoadProperties(defaultPrefix, DBProperty{}).(DBProperty)
	// 连接池信息
	c.defaultPoolProperty = container.LoadProperties(defaultPrefix, PoolProperty{}).(PoolProperty)
	// 日志信息
	c.defaultLogProperty = container.LoadProperties(defaultPrefix, LogProperty{}).(LogProperty)

	var tags []string
	if val := container.GetProperty("gorm.multi"); val != nil {
		tags = strings.Split(val.(string), ",")
	}
	if len(tags) == 0 {
		db := c.generateDB(c.defaultDBProperty, c.defaultPoolProperty, c.defaultLogProperty)
		// 注册db
		baseDao := dao.New(db)
		container.RegisterNamedBean("gorm", db)
		container.RegisterBean(baseDao)
		c.log.Info(container.Context(), "Gorm library loaded")
		return
	}
	for _, tag := range tags {
		prefix := defaultPrefix + tag + "."
		// db连接信息
		dbProperty := container.LoadProperties(prefix, DBProperty{}).(DBProperty)
		dbProperty.Merge(c.defaultDBProperty)
		// 连接池信息
		poolProperty := container.LoadProperties(prefix, PoolProperty{}).(PoolProperty)
		poolProperty.Merge(c.defaultPoolProperty)
		// 日志信息
		logProperty := container.LoadProperties(prefix, LogProperty{}).(LogProperty)
		logProperty.Merge(c.defaultLogProperty)

		db := c.generateDB(dbProperty, poolProperty, logProperty)
		// 注册db
		baseDao := dao.New(db)
		container.RegisterNamedBean("gormFor"+tag, db)
		container.RegisterNamedBean(tag, baseDao)

		c.log.Info(container.Context(), "Gorm library loaded db "+tag)
	}
}

func (c *configuration) generateDB(dbProperty DBProperty, pool PoolProperty, logProperty LogProperty) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", []interface{}{
		dbProperty.Username,
		dbProperty.Password,
		dbProperty.Host,
		dbProperty.Port,
		dbProperty.Database,
		dbProperty.Parameters,
	}...)
	// 配置数据库
	db, err := gorm.Open(mysql.Open(dsn), c.opts...)
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	// 配置连接池
	if pool.MaxLifeTime > 0 {
		sqlDB.SetConnMaxLifetime(pool.MaxLifeTime)
	}
	if pool.MaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(pool.MaxIdleTime)
	}
	sqlDB.SetMaxIdleConns(pool.MaxIdle)
	sqlDB.SetMaxOpenConns(pool.MaxOpen)
	// 配置日志
	db.Logger = newLogger(c.log, logProperty)
	return db
}
