---
sidebar_label: 'redis-write'
title: redis-write
sidebar_position: 5
---
## redis读写

#### 版本记录

日期 |内容|撰写人
|---|---|---|
2023/02/09| V1.0 初版| 刘胜杰

## 人员
- 需求提出：刘胜杰
- 研发： 刘胜杰
- 抄送：

## 问题

支持redis的数据查询、命令执行

## 目标

- 支持数据查询
- 支持写命令审批、执行
- 读写命令，支持白名单、影响范围限制等

## 设计
### 整体设计

- 读数据
    - 限制范围内，直接执行读取并返回结果
- 写数据
    - 生成任务、通过审批后可执行
    - 一个任务内科包含多条命令
- 集群管理复用DB集群管理，增加Type字段，支持UI检索


### 详细设计

- 在白名单的基础上，支持范围限制
- 读白名单

    命令|说明|范围
    |---|---|---|
    get |
    mget |
    hget |
    hmget |
    lrange |限制影响数据量|1000
    zrange |限制影响数据量|1000
    sismember ||
    scard ||
    zcard ||
    hscan ||
    ttl ||
    type ||
    hlen ||
    exists ||
    sscan ||

- 写白名单

    命令|说明|范围
    |---|---|---|
    set ||
    mset ||
    hset ||
    hmset ||
    hdel ||
    del ||
    zrem ||
    srem ||
    sadd ||
    zadd ||
    incrby ||
- ……

排期

- [x] 暂无[人员][日期] 
- [ ] Milestone 2 
- [ ] ……