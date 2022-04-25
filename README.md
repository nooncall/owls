
<div align=center>
<img src="https://img0.baidu.com/it/u=2822765666,2555722031&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=501" width=300" height="300" />
</div>
<div align=center>
<img src="https://img.shields.io/badge/golang-1.16-blue"/>
<img src="https://img.shields.io/badge/gin-1.7.0-lightBlue"/>
<img src="https://img.shields.io/badge/vue-3.2.25-brightgreen"/>
<img src="https://img.shields.io/badge/element--plus-2.0.1-green"/>
<img src="https://img.shields.io/badge/gorm-1.22.5-red"/>
</div>

[English](./README-en.md) | 简体中文

## 版本

初版开发中，目前tidb、mysql的sql审核功能可用。 


## 0. 基本介绍

### 0.1 项目介绍

> Owls是一个基于 [vue](https://vuejs.org) 和 [gin](https://gin-gonic.com) 开发的全栈前后端分离的后端系统管理平台，集成jwt鉴权，sql审批（sql查询，mq管理、redis使用管理、etcd管理等）。帮助您更方便、更规范的管理后端中间件系统，守护数据系统，提高系统稳定性，避免人为操作失误导致的故障。

[在线预览](http://106.13.50.14/owls): http://106.13.50.14/owls

测试用户名：admin

测试密码：123456

### 0.2 贡献指南
Hi! 首先感谢你使用 Owls。

Owls 是一套为快安全、规范准备的一整套前后端分离架构式的后端系统管理平台, 旨在简化、标准化DBA、研发等同学对数据库、MQ、缓存等系统的使用和管理。

Owls 的成长离不开大家的支持，如果你愿意为 Owls 贡献代码或提供建议，请阅读以下内容。

#### 0.2.1 Issue 规范
- 稍大的功能模块设计需要在文档目录（server/docs/design）中，参照[模板](./server/docs/design/template.md)添加设计说明，并进行review。

- issue 仅用于提交 Bug 或 Feature 以及设计相关的内容。

- 在提交 issue 之前，请搜索相关内容是否已被提出。

#### 0.2.2 Pull Request 规范
- 请先 fork 一份到自己的项目下，不要直接在仓库下建分支。

- commit 信息要以`[模块名]: 描述信息` 的形式填写，例如 `[readme] update xxx msg`。

- 如果是修复 bug，请在 PR 中给出描述信息。

- 合并代码需要一人进行 review 后 approve，方可合并。

## 1. 开发计划

* db查询支持
* db分库分表update支持
* mq使用管理设计及实现
* redis使用管理设计及实现
* etcd使用管理设计及实现
* db代理设计及实现
* 缓存代理设计及实现

## 2. 使用说明

```
- node版本 > v12.18.3
- golang版本 >= v1.16
- IDE推荐：Goland
- 初始化项目： mysql数据库
- 替换掉项目中的七牛云公钥，私钥，仓名和默认url地址，以免发生测试文件数据错乱
```

### 2.1 server项目

使用 `Goland` 等编辑工具，打开server目录，不可以打开 Owls 根目录

```bash

# 克隆项目
git clone https://github.com/qingfeng777/owls.git
# 进入server文件夹
cd server

# 使用 go mod 并安装go依赖包
go mod download

# 编译 
go build -o server main.go (windows编译命令为go build -o server.exe main.go )

# 运行二进制
./server (windows运行命令为 server.exe)
```

### 2.2 web项目

```bash
# 进入web文件夹
cd web

# 安装依赖
cnpm install || npm install

# 启动web项目
npm run serve
```

### 2.3 swagger自动化API文档

#### 2.3.1 安装 swagger

##### （1）可以访问外国网站

````
go get -u github.com/swaggo/swag/cmd/swag
````

##### （2）无法访问外国网站

由于国内没法安装 go.org/x 包下面的东西，推荐使用 [goproxy.cn](https://goproxy.cn) 或者 [goproxy.io](https://goproxy.io/zh/)

```bash
# 如果您使用的 Go 版本是 1.13 - 1.15 需要手动设置GO111MODULE=on, 开启方式如下命令, 如果你的 Go 版本 是 1.16 ~ 最新版 可以忽略以下步骤一
# 步骤一、启用 Go Modules 功能
go env -w GO111MODULE=on 
# 步骤二、配置 GOPROXY 环境变量
go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,direct

# 如果嫌弃麻烦,可以使用go generate 编译前自动执行代码, 不过这个不能使用 `Goland` 或者 `Vscode` 的 命令行终端
cd server
go generate -run "go env -w .*?"

# 使用如下命令下载swag
go get -u github.com/swaggo/swag/cmd/swag
```

#### 2.3.2 生成API文档

```` shell
cd server
swag init
````

> 执行上面的命令后，server目录下会出现docs文件夹里的 `docs.go`, `swagger.json`, `swagger.yaml` 三个文件更新，启动go服务之后, 在浏览器输入 [http://localhost:8888/swagger/index.html](http://localhost:8888/swagger/index.html) 即可查看swagger文档


## 3. 技术选型

- 前端：用基于 [Vue](https://vuejs.org) 的 [Element](https://github.com/ElemeFE/element) 构建基础页面。
- 后端：用 [Gin](https://gin-gonic.com/) 快速搭建基础restful风格API，[Gin](https://gin-gonic.com/) 是一个go语言编写的Web框架。
- 数据库：采用`MySql`(5.6.44)版本，使用 [gorm](http://gorm.cn) 实现对数据库的基本操作。
- 缓存：使用`Redis`实现记录当前活跃用户的`jwt`令牌并实现多点登录限制。
- API文档：使用`Swagger`构建自动化文档。
- 配置文件：使用 [fsnotify](https://github.com/fsnotify/fsnotify) 和 [viper](https://github.com/spf13/viper) 实现`yaml`格式的配置文件。
- 日志：使用 [zap](https://github.com/uber-go/zap) 实现日志记录。

## 4. 项目架构

### 4.1 系统架构图

![系统架构图](http://qmplusimg.henrongyi.top/gva/gin-vue-admin.png)


## 5. 主要功能
### Sql审核
- sql任务流审批。
- sql执行及备份、回滚。
- 规范化规则自动审批。
- 集群管理。

### 基础功能：
- 权限管理：基于`jwt`和`casbin`实现的权限管理。
- 用户管理：系统管理员分配用户角色和角色权限。
- 角色管理：创建权限控制的主要对象，可以给角色分配不同api权限和菜单权限。
- 菜单管理：实现用户动态菜单配置，实现不同角色不同菜单。
- 条件搜索：增加条件搜索示例。
- 多点登录限制：需要在`config.yaml`中把`system`中的`use-multipoint`修改为true(需要自行配置Redis和Config中的Redis参数，测试阶段，有bug请及时反馈)。

### [关于我们]()

## 6. 鸣谢

### 感谢 [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin) 项目，是一个很好用的基础脚手架。
### 感谢 [伴鱼Owl](https://github.com/ibanyu/owl)项目，此项目继承于伴鱼Owl。