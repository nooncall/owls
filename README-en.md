
# Owls - Data Interaction Management Platform

<div align="center">
<img src="https://img0.baidu.com/it/u=2822765666,2555722031&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=501" width="300" height="300" />
</div>
<div align="center">
<img src="https://img.shields.io/badge/golang-1.16-blue"/>
<img src="https://img.shields.io/badge/gin-1.7.0-lightBlue"/>
<img src="https://img.shields.io/badge/vue-3.2.25-brightgreen"/>
<img src="https://img.shields.io/badge/element--plus-2.0.1-green"/>
<img src="https://img.shields.io/badge/gorm-1.22.5-red"/>
</div>

English | [简体中文](./README.md)

## Introduction

Owls is a full-stack, front-end and back-end separated data interaction management platform developed with [Vue.js](https://vuejs.org) and [Go](https://go.dev/). It integrates JWT authentication and provides features such as SQL approval (SQL query, MQ management, Redis usage management, etcd management, etc.). It helps you manage middleware systems more conveniently and effectively, safeguard data systems, improve system stability, prevent faults caused by human errors, enhance efficiency, and unleash creativity.

### Project Links

Github (access with VPN): https://github.com/nooncall/owls  
Gitee (access from China): https://gitee.com/nooncall/owls

### Online Preview

[Online Preview](http://owls.nooncall.cn:8778/owls): http://owls.nooncall.cn:8778/owls

Test User: admin  
Password: aaaaaa

### Features

#### SQL Approval

- Data querying
- SQL task approval
- SQL execution, backup, and rollback
- Automated approval based on standardization rules
- Cluster management

#### Basic Features

- Access control: Role-based access control (RBAC) implemented with JWT and Casbin
- User management: System administrators can assign user roles and role permissions
- Role management: Create main objects for access control, assign different API permissions and menu permissions to roles
- Conditional search: Add an example of conditional search
- Support for registration and LDAP login methods

### Documentation

#### Quick Start

Please refer to the [Quick Start](./docs/user_guide/quick-start.md) document or visit owls.nooncall.cn:8778/docs/user_guide/quick-start.

#### Complete Documentation

Please visit owls.nooncall.cn:8778/docs.

After the project is started, the documentation accompanying the deployment can also be accessed at the project's addr:port/docs directory.

#### Contribution Guidelines

Please refer to the [Contribution Guidelines](./docs/roadmap/contribution.md) document or visit owls.nooncall.cn/docs/roadmap/contribution.

### Acknowledgments

- Special thanks to the [Owl](https://github.com/ibanyu/owl) project by ibanyu, which Owls inherits from.
- Many thanks to the [Docusaurus](https://github.com/facebook/docusaurus) project, a very useful documentation system.
- Thanks to the [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin) project, a powerful and practical scaffold.