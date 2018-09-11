
---
## 如何建立集群&动态增删节点

> 如何使用该服务

#### 一、 我们需要安装etcdv3服务
##### 1.download etcdv3二进制文件

> https://github.com/coreos/etcd/releases/		（etcd-v3.3.9-linux-amd64.tar.gz）

##### 2.解压得到两个文件
> - etcd
> - etcdctl

##### 3. 复制到 /user/local/bin 目录下

##### 4. 配置环境变量

> ETCDCTL_API=3

##### 5.在命令行测试 etcd 
> etcd

##### 6.在命令行测试 etcdctl
> etcdctl

#### 二、 我们需要部署etcdv3服务集群
##### 1.通用配置信息（假设实现三台电脑的etcdv3集群）
###### 每台电脑都需要配置

```
TOKEN=token-07
CLUSTER_STATE=new
NAME_1=machine-1
NAME_2=machine-2
NAME_3=machine-3
HOST_1=10.204.29.77
HOST_2=10.204.29.70
HOST_3=10.204.29.73
CLUSTER=${NAME_1}=http://${HOST_1}:2380,${NAME_2}=http://${HOST_2}:2380,${NAME_3}=http://${HOST_3}:2380
```

##### 2.为各电脑启动 etcdv3
###### 电脑一（IP:10.204.29.77）

```
THIS_NAME=${NAME_1}
THIS_IP=${HOST_1}
etcd --data-dir=data.etcd --name ${THIS_NAME} --initial-advertise-peer-urls http://${THIS_IP}:2380 --listen-peer-urls http://${THIS_IP}:2380 --advertise-client-urls http://${THIS_IP}:2379 --listen-client-urls http://${THIS_IP}:2379 --initial-cluster ${CLUSTER} --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
```

###### 电脑二（IP:10.204.29.70）


```
THIS_NAME=${NAME_2}
THIS_IP=${HOST_2}
etcd --data-dir=data.etcd --name ${THIS_NAME} --initial-advertise-peer-urls http://${THIS_IP}:2380 --listen-peer-urls http://${THIS_IP}:2380 --advertise-client-urls http://${THIS_IP}:2379 --listen-client-urls http://${THIS_IP}:2379 --initial-cluster ${CLUSTER} --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
```


###### 电脑三（IP:10.204.29.73）



```
THIS_NAME=${NAME_3}
THIS_IP=${HOST_3}
etcd --data-dir=data.etcd --name ${THIS_NAME} --initial-advertise-peer-urls http://${THIS_IP}:2380 --listen-peer-urls http://${THIS_IP}:2380 --advertise-client-urls http://${THIS_IP}:2379 --listen-client-urls http://${THIS_IP}:2379 --initial-cluster ${CLUSTER} --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
```

##### 3.查看集群列表（任意一台主机都可以）出现集群列表表示服务成功启动


```
HOST_1=10.204.29.77
HOST_2=10.204.29.70
HOST_3=10.204.29.73
ENDPOINTS=$HOST_1:2379,$HOST_2:2379,$HOST_3:2379

etcdctl --endpoints=$ENDPOINTS member list
```



##### 4.以表格形式查看&&健康状况反馈

> etcdctl --write-out=table --endpoints=$ENDPOINTS endpoint status


```
+-------------------+------------------+---------+---------+-----------+-----------+------------+
|     ENDPOINT      |        ID        | VERSION | DB SIZE | IS LEADER | RAFT TERM | RAFT INDEX |
+-------------------+------------------+---------+---------+-----------+-----------+------------+
| 10.204.29.77:2379 | 3219c8c0de0d6d77 |   3.3.9 |   20 kB |      true |        14 |          9 |
| 10.204.29.70:2379 | ac0a30888be2343b |   3.3.9 |   20 kB |     false |        14 |          9 |
| 10.204.29.73:2379 | ef3b7488a017dabf |   3.3.9 |   20 kB |     false |        14 |          9 |
+-------------------+------------------+---------+---------+-----------+-----------+------------+
```


> etcdctl --endpoints=$ENDPOINTS endpoint health


```
10.204.29.77:2379 is healthy: successfully committed proposal: took = 3.539859ms
10.204.29.70:2379 is healthy: successfully committed proposal: took = 3.177144ms
10.204.29.73:2379 is healthy: successfully committed proposal: took = 4.091101ms
```



##### 5.删除某个节点

> ###### 在leader上设置要删除节点的ID



```
MEMBER_ID=ef3b7488a017dabf
```


> ###### 执行删除


```
etcdctl --endpoints=${HOST_1}:2379,${HOST_2}:2379,${HOST_3}:2379 member remove ${MEMBER_ID}
```

> - ###### 预期结果

>  etcdctl --write-out=table --endpoints=$ENDPOINTS endpoint status

>  etcdctl --endpoints=$ENDPOINTS endpoint health



```
+-------------------+------------------+---------+---------+-----------+-----------+------------+
|     ENDPOINT      |        ID        | VERSION | DB SIZE | IS LEADER | RAFT TERM | RAFT INDEX |
+-------------------+------------------+---------+---------+-----------+-----------+------------+
| 10.204.29.77:2379 | 3219c8c0de0d6d77 |   3.3.9 |   20 kB |      true |        14 |         10 |
| 10.204.29.70:2379 | ac0a30888be2343b |   3.3.9 |   20 kB |     false |        14 |         10 |
+-------------------+------------------+---------+---------+-----------+-----------+------------+
```



```
10.204.29.77:2379 is healthy: successfully committed proposal: took = 2.646736ms
10.204.29.70:2379 is healthy: successfully committed proposal: took = 3.52799ms
10.204.29.73:2379 is unhealthy: failed to connect: dial tcp 10.204.29.73:2379: connect: connection refused
Error: unhealthy cluster
```

###### 电脑三（IP:10.204.29.73）


```
2018-09-02 22:30:52.085144 I | rafthttp: stopped peer ac0a30888be2343b
```


##### 6.添加某节点

###### 电脑四 （IP: 10.204.30.188）
- 假设节点三已经被删除


```
export ETCDCTL_API=3
NAME_1=machine-1
NAME_2=machine-2
NAME_4=machine-4
HOST_1=10.204.29.77
HOST_2=10.204.29.70
HOST_4=10.204.30.188
etcdctl --endpoints=${HOST_1}:2379,${HOST_2}:2379 member add ${NAME_4} --peer-urls=http://${HOST_4}:2380
```



```
TOKEN=token-07
CLUSTER_STATE=existing
NAME_1=machine-1
NAME_2=machine-2
NAME_4=machine-4
HOST_1=10.204.29.77
HOST_2=10.204.29.70
HOST_4=10.204.30.188
CLUSTER=${NAME_1}=http://${HOST_1}:2380,${NAME_2}=http://${HOST_2}:2380,${NAME_4}=http://${HOST_4}:2380
```



```
THIS_NAME=${NAME_4}
THIS_IP=${HOST_4}

etcd --data-dir=data.etcd --name ${THIS_NAME} --initial-advertise-peer-urls http://${THIS_IP}:2380 --listen-peer-urls http://${THIS_IP}:2380 --advertise-client-urls http://${THIS_IP}:2379 --listen-client-urls http://${THIS_IP}:2379 --initial-cluster ${CLUSTER} --initial-cluster-state ${CLUSTER_STATE} --initial-cluster-token ${TOKEN}
```


###### 此时节点成功接入

> ###### 在电脑一（IP:10.204.29.77）进行验证
 

###### 设置集群信息

```
HOST_1=10.204.29.77
HOST_2=10.204.29.70
HOST_4=10.204.30.188
ENDPOINTS=$HOST_1:2379,$HOST_2:2379,$HOST_4:2379
```

###### 查看状态


```
etcdctl --write-out=table --endpoints=$ENDPOINTS endpoint status
```


```
etcdctl --endpoints=$ENDPOINTS endpoint health
```



```
+--------------------+------------------+---------+---------+-----------+-----------+------------+
|      ENDPOINT      |        ID        | VERSION | DB SIZE | IS LEADER | RAFT TERM | RAFT INDEX |
+--------------------+------------------+---------+---------+-----------+-----------+------------+
|  10.204.29.77:2379 | 3219c8c0de0d6d77 |   3.3.9 |   20 kB |      true |        16 |         12 |
|  10.204.29.70:2379 | ac0a30888be2343b |   3.3.9 |   20 kB |     false |        16 |         12 |
| 10.204.30.188:2379 |  10e1765611b5901 |   3.3.9 |   20 kB |     false |        16 |         12 |
+--------------------+------------------+---------+---------+-----------+-----------+------------+
```




```
10.204.30.188:2379 is healthy: successfully committed proposal: took = 3.585037ms
10.204.29.70:2379 is healthy: successfully committed proposal: took = 4.057687ms
10.204.29.77:2379 is healthy: successfully committed proposal: took = 3.077686ms
```



##### 本内容不涉及代码 详细内容参照官方手册
##### URL：[https://github.com/etcd-io/etcd/blob/master/Documentation/demo.md](https://github.com/etcd-io/etcd/blob/master/Documentation/demo.md)
###### 编著人：Allen guo
###### 日期：  2018/9/3 




























