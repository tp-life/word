### 娱堂平台

#### 开发说明

- 路由建议是我们自己后端接口前缀加上 /admin, 开放平台的接口前缀加上 /open, 其他接口加上 /api
- 使用 gin 的三种mode进行配置环境的管理, 分别为 debug release 和 test
- 所有代码尽量符合 golint 规则
- 尽量避免使用 interface{} 对象, 所有数据尽量结构化, 包括http返回数据和表单验证等, 而不是直接使用 map[string]interface{}, 
如果存在需要接受多种类型的方法, 请尽量使用接口, 在实在无法解决时采用 interface{}
- 注释完备
- 虽然产品设计为好几端, 比如开发者平台, 官网等, 但是后端开发按照各个模块统一开发, 但是要保证所有需求能够提供确切功能,
如果在正式服部署时需要不同服务不同路由, 由 nginx 来转发, 但是在自己cmd/xxx下注册服务时, 记得分组
- 每个模块应该有一份 README.md 说明自己
#### 功能

##### 目录结构
```code 
├── Makefile
├── README.md
├── api
├── assets
├── build
├── cmd
├── configs
│   ├── README.md
│   ├── debug
│   ├── release
│   └── test
├── docs
├── examples
├── go.mod
├── go.sum
├── internal
├── locales
│   ├── README.md
│   ├── en.toml
│   └── zh.toml
├── logs
├── pkg
│   ├── app
│   ├── captcha
│   ├── database
│   ├── elastic
│   ├── email
│   ├── i18n
│   ├── log
│   ├── middlewares
│   ├── pager
│   ├── password
│   ├── payments
│   ├── queue
│   ├── rbac
│   ├── redis
│   ├── sensitivewords
│   ├── server
│   ├── sessions
│   ├── unique
│   └── validator
├── scripts
├── test
├── tools
└── web

```
### 你首先需要了解的目录
##### cmd 
cmd 程序 main 入口

##### configs
configs 是配置目录

##### internal 
程序业务代码,注意分层
* internal 中第一层应该为各个服务, 比如 user, payment, order, news 等
* 内层以 user 举例, 建议为 user->entity, user->service, user->controller 等
* entity 层主要是数据表结构和自定义的一些操作
* service 是业务主要逻辑, 这里不要包含 gin 的 Context , 由 controller 统一调用

##### locales
国际化翻译文件

##### logs 
日志文件

##### pkg 
公共调用模块, 这里放置所有模块都能使用的代码, 但不包括 pkg 本身, 因为会造成循环依赖的问题

##### 编译
```bash
# Windows
make o=yutang.exe mode=debug|test|release build

#Linux
make o=yutang mode=debug|test|release build
```