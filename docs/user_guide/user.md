---
sidebar_label: '用户'
title: 用户管理
sidebar_position: 2
---

### 初始化用户

执行初始化时，会初始化一个`admin`用户，一个`user`用户，其密码是在配置文件server/config.yaml中配置的，默认为`aaaaaa`

### 新用户接入方式

新用户接入支持两种方式，由server/config.yaml配置文件中的login/model控制，registry为注册登录模式，ldap为接入LDAP的方式。默认值为registry。

#### 注册-登陆

配置为注册登录模式，会在用户页面显示`注册`按钮,点击即可进行注册，新注册的用户默认为普通用户权限，如需调整调整权限，可由admin用户创建并赋予其新的角色。

#### LDAP接入用户

LDAP模式需要正确填写server/config.yaml配置文件中login/ldap下的配置。配置正确后可接入LDAP系统，使用LDAP账号密码登陆系统。新登陆的用户默认同样为普通用户权限。