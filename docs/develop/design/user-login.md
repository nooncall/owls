---
sidebar_label: '登录'
title: 登录
sidebar_position: 2
---
## 标题

#### 版本记录

日期 |内容|撰写人
|---|---|---|
2022/06/05| V1.0 初版| 刘胜杰

## 人员
- 需求提出：
- 研发： 刘胜杰 & 
- 抄送：

## 问题
  解决不用场景下的用户接入，用户登录问题

## 目标
  提供多种方式，支持不同场景下用户接入需求  
  1，支持ldap接入方式  
  2，支持注册登录方式  
  3，设计好用户接口，使得方便接入企业内部用户系统

## 设计
- 以配置切换不同的用户模式  
    根据配置值，注入不同的用户登录实现，
    并根据配置值展示登陆页面的是否有‘注册’按钮
- 梳理现有的用户系统，设计登陆相关的接口  
    需改造现有的登陆实现，
- 支持ldap接入用户登录  
    现有部分ldap代码，需要接入改造并测试
- 支持注册登录方式  
    登陆已经存在，需添加注册功能  
    目前管理员可以修改用户密码，需支持用户修改自己的密码

## 排期

- [x] 支持ldap  2022-06中 
- [x] 支持注册登录方式   2022-06下
- [ ] 整理接入内部用户系统的文档