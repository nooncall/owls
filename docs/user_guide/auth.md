---
sidebar_label: '权限'
sidebar_position: 3
---

### 权限实现

权限控制采用[Casbin](https://casbin.org/docs/zh-CN/overview)框架，使用的RBAC模型，默认会初始化一个admin角色和一个user角色，分别对应于[用户](./%E7%94%A8%E6%88%B7)中的admin用户和普通用户。

如需创建新的角色，可直接在`超级管理员/角色管理`页面创建新的角色并勾选赋予其权限即可。