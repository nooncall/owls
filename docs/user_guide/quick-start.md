---
sidebar_label: '快速开始'
sidebar_position: 1
---

### Docker安装

如果使用准备好的DB，则可以直接执行：

    docker run -p 8081:80 -d  mingbai/owls:latest

如无，则创建docker网桥，然后分别docker启动mysql、owls的容器：

    docker network create owls

    docker run -d --network=owls --name=mysql  -e MYSQL_ROOT_PASSWORD=aaaaaa mysql:5.7

    docker run -p 8081:80 -d  --network=owls  mingbai/owls:latest

### 集群内安装

    kubectl apply -n argo -f https://github.com/nooncall/owls/tree/master/docs/user_guide/deployment.yaml


### 初始化

登陆页面点击初始化按钮，根据我们上面的安装步骤，数据库的地址应该写`mysql`(或自行准备的DB地址),密码`aaaaaa`,其他默认即可。

### 本地文档

访问： `http://localhost:8081/docs`

### Enjoy
现在访问`http://localhost:8081`即可尝试使用系统提供的功能了 。 

默认创建的用户有两个，`admin`和`user` 密码都是`aaaaaa` 。

正式使用的系统，还是建议使用独立的、持久化的数据库，可以选择Mysql或者TiDB。