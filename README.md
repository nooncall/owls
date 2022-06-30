
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

### 介绍

> Owls是一个基于 [vue](https://vuejs.org) 和 [go](https://go.dev/) 开发的全栈前后端分离的数据交互管理平台，集成jwt鉴权，sql审批（sql查询，mq管理、redis使用管理、etcd管理等）。帮助您更方便、更规范的管理中间件系统，守护数据系统，提高系统稳定性，避免人为操作失误导致的故障，提高效率，解放创造力。

### 在线预览

[在线预览](http://106.13.50.14/owls): http://106.13.50.14/owls

测试用户名：admin

密码：aaaaaa

### 功能介绍

#### Sql审核功能

- 数据查询
- sql任务流审批。
- sql执行及备份、回滚。
- 规范化规则自动审批。
- 集群管理。

### 基础功能

- 权限管理：基于`jwt`和`casbin`实现的权限管理。
- 用户管理：系统管理员分配用户角色和角色权限。
- 角色管理：创建权限控制的主要对象，可以给角色分配不同api权限和菜单权限。
- 条件搜索：增加条件搜索示例。
- 支持注册、接入LDAP两种登陆方式。

### 文档

#### 快速开始

请查询[快速开始](./docs/user_guide/quick-start.md) 

或者访问 owls.nooncall.cn/docs/user_guide/quick-start

#### 完整文档

请查阅 owls.nooncall.cn 。  

项目启动后，在项目addr:port/docs 目录下也能查阅跟随项目部署的文档。

#### 贡献指南

请查询[贡献指南](./docs/roadmap/contribution.md)  

或者访问 owls.nooncall.cn/docs/roadmap/contribution

### 鸣谢

##### 感谢 [伴鱼Owl](https://github.com/ibanyu/owl)项目，此项目继承于伴鱼Owl。
##### 感谢[docusaurus](https://github.com/facebook/docusaurus)项目，一个非常好用的文档系统。
##### 感谢 [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin) 项目，是一个非常好用的基础脚手架。