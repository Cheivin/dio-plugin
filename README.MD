# DIO插件

## gin插件

### 使用

- 安装 `go get github.com/cheivin/dio-plugin/gin`

- 配置 `dio.use(web.GinWeb(useLogger, useCors))`

- Dio中注入

  ```go
  package main
  import (
   "github.com/gin-gonic/gin"
  )
  
  type Api struct{
      Web *gin.Engine `aware:"web"`
  }
  ```

### 配置参数

| 配置项                            | 默认值          | 说明                                        |
|--------------------------------|--------------|-------------------------------------------|
| app.port                       | 8080         | 启动端口                                      |
| app.web.log.skip-path          |              | 跳过日志的路径。useLogger为true时生效                 |
| app.web.log.trace-name         | X-Request-Id | 会话跟踪ID,定义在header中。useLogger为true时生效       |
| app.web.cors.origin            |              | 允许跨域的域名,多个用`,`分割，为空时表示不限制。useCors为true时生效 |
| app.web.cors.method            |              | 允许跨域的方法,多个用`,`分割，为空时表示不限制。useCors为true时生效 |
| app.web.cors.header            |              | 允许跨域的header,多个用`,`分割。useCors为true时生效      |
| app.web.cors.allow-credentials | true         | 是否允许发送cookie。useCors为true时生效              |
| app.web.cors.expose-header     |              | 暴露的header,多个用`,`分割。useCors为true时生效        |
| app.web.cors.max-age           | 43200        | 预检有效期,单位秒。useCors为true时生效                 |

## gorm

### 使用

- 安装 `go get github.com/cheivin/dio-plugin/gorm`

- 配置 `dio.use(web.Gorm(options))`

- Dio中注入

  ```go
  package main
  import (
   "github.com/cheivin/dio/plugin/gorm/dao"
  )
  
  type XXDao struct{
     *dao.Dao `aware:""`
  }
  ```

### 配置参数

| 配置项                      | 默认值       | 说明                                                                                       |
|--------------------------|-----------|------------------------------------------------------------------------------------------|
| gorm.username            | root      | 用户名                                                                                      |
| gorm.password            | root      | 密码                                                                                       |
| gorm.host                | localhost | 主机                                                                                       |
| gorm.port                | 3306      | 端口                                                                                       |
| gorm.database            |           | 数据库                                                                                      |
| gorm.parameters          |           | 附加参数.[dsn-data-source-name](https://github.com/go-sql-driver/mysql#dsn-data-source-name) |
| gorm.pool.max-idle       | 0         | 最大空闲连接数,0表示不限制                                                                           |
| gorm.pool.max-open       | 0         | 最大连接数,0表示不限制                                                                             |
| gorm.pool.max-life-time  |           | 连接最大重用时间,duration格式                                                                      |
| gorm.pool.max-idle-time  |           | 连接最大空闲时间,duration格式                                                                      |
| gorm.log.level           | 4         | 日志级别,1:静默;2:错误;3:警告,4:普通                                                                 |
| gorm.log.slow-log        |           | 慢日志阈值,duration格式。为空则不开启                                                                  |
| gorm.log.ignore-notfound |           | 是否忽略notfound类型错误日志                                                                       |
