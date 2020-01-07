

# zerodb 数据源切换

## zerodb 集群

* 3台proxy(10.1.21.82 、10.1.21.79、10.1.21.242) 、1台keeper(10.1.24.173)、 一个etcd集群(10.1.24.173)、
* proxy 集群cluster-func-test
* 8台mysql 4个hostgroups(10.1.22.1-4)(10.1.21.79-82)

![集群架构](https://processon.com/chart_image/5a77c423e4b064e9ddb7172b.png)

## 数据源切换配置


* NeedVote   //切换时是否需要投票,true 允许、false 不允许、默认true
* VoteApproveRatio  //投票比列, 超过这个比例即认为集群中所有zeroproxy允许切换、默认50%
* NeedLoadCheck  //是否需要检查mysql load,true 允许、false 不允许、默认true
* SafeLoad  //安全的load值,超过改值不允许切换 、默认4
* NeedBinlogCheck  //是否需要检查binlog、默认true
* SafeBinlogDelay  //安全的binlog 延迟数值. 默认1000


## Proxy 下线恢复

* OfflineOnLostKeeper 	  //当proxy被拒绝多次后是否向keeper注销自己
* OfflineSwhRejectedNum.  //proxy被拒绝几次后加入注销队列
* OfflineDownHostNum       //跟多少台mysql host断开后向注销自己
* OfflineRecover.                  //跟多台mysql host断开后是否下线自己默认必须为true

## 参数


* Hostgroup    // hostgroup Name
* From.            // 需要切换的mysql host索引 (当前连接的) mysql host 索引
* To.                // 切换到的mysql host 索引
* ProxyIP        // 发起切换的ZeroProxy IP.
* ClusterName  //ZeroProxy 的集群名字


## 数据库切换流程图

![](http://on-img.com/chart_image/5b7d5b9ce4b06fc64ad0b53f.png?_=1534990138970)




## Case1. Mysql 产生问题

### 1. mysql进程挂掉,所有proxy无法ping通mysql

**预期结果**:keeper返回无法读取binlog日志,拒绝切换.

1. 查看当前数据源

```
route select * from order_ins;
```

2. 停止hostgroup3 的10.1.22.3 mysql

```
ansible -i ./zerodb-ansible/hosts 10.1.22.3 -m "shell" -a "service mysqld stop " --extra-vars "ansible_user=root ansible_password=sdfsdf"
```

查看keeper日志

```
ansible -i ./zerodb-ansible/hosts keeper_daily -m "shell" -a "ls /var/log/zlog/" --extra-vars "ansible_user=root ansible_password=sdfsdf"

ansible -i ./zerodb-ansible/hosts keeper_daily -m "shell" -a "tail -200 /var/log/zlog/ds_switch.log.2018083000" --extra-vars "ansible_user=root ansible_password=sdfsdf"
```

查看proxy日志

```
ansible -i ./zerodb-ansible/hosts proxy_daily -m "shell" -a "tail -100 /var/log/zlog/ds_switch.log.2018083000" --extra-vars "ansible_user=root ansible_password=sdfsdf" 
```



### 2. Mysql 网络问题,所有proxy无法ping通mysql

 **预期结果**: keeper根据投票结果,发生切换.

1. 使用iptables 禁止所有proxy访问hostgroup 10.1.22.3

```
service mysqld start
iptables -I INPUT -p tcp -s 10.1.21.79 --dport 3306 -j DROP
iptables -I INPUT -p tcp -s 10.1.21.82 --dport 3306 -j DROP
iptables -I INPUT -p tcp -s 10.1.21.242 --dport 3306 -j DROP
iptables -L INPUT --line-numbers
iptables -D INPUT 1
```

2. 查看keepr切换日志


```
ansible -i ./zerodb-ansible/hosts keeper_daily -m "shell" -a "tail -200 /var/log/zlog/ds_switch.log.2018083000" --extra-vars "ansible_user=root ansible_password=sdfsdf"
```

3. 查看proxy切换日志

```
ansible -i ./zerodb-ansible/hosts proxy_daily -m "shell" -a "tail -100 /var/log/zlog/ds_switch.log.2018083000" --extra-vars "ansible_user=root ansible_password=sdfsdf" 
```

4. 当前连接的mysql

```
route select * from order_ins;
```

5. 设置iptables 使所有proxy无法访问10.1.21.81.数据库切回10.1.22.3


```
iptables -D INPUT 1 //删除10.1.22.3 iptables
iptables -L INPUT --line-numbers
iptables -I INPUT -p tcp -s 10.1.21.79 --dport 3306 -j DROP
iptables -I INPUT -p tcp -s 10.1.21.82 --dport 3306 -j DROP
iptables -I INPUT -p tcp -s 10.1.21.242 --dport 3306 -j DROP
```


### 3. mysql负载过高. 

**预期结果**:Keeper 拒绝切换.

1. 配置SafeLoad 为2

```
etcdctl --endpoints=10.1.24.173:2379 get --prefix zerodb/shardconf/cluster-func-test/default/switch 
```

2. 设置iptables 使数据库从10.1.22.3 切到10.1.21.81. 

3. 使用stress 让 10.1.21.81 load变高

```
stress --cpu 1 --timeout 180
```

4. 查看切换日志 

```
ansible -i ./zerodb-ansible/hosts keeper_daily -m "shell" -a "tail -200 /var/log/zlog/ds_switch.log.2018082900" --extra-vars "ansible_user=root ansible_password=sdfsdf"

ansible -i ./zerodb-ansible/hosts proxy_daily -m "shell" -a "tail -100 /var/log/zlog/ds_switch.log.2018082900" --extra-vars "ansible_user=root ansible_password=sdfsdf" 
```





## Case2. 单台proxy网络产生问题

### 1. 无法连接一台mysql

**预期结果**: keeper拒绝切换,proxy 无法获得足够的投票.

1.  proxy 10.1.21.82设置iptables 无法访问10.1.22.3 

```
etcdctl --endpoints=10.1.24.173:2379 get --prefix zerodb/proxy/status/cluster-func-test/
```

```
iptables -A OUTPUT  -p tcp  -d 10.1.22.3  --dport 3306 -j REJECT
iptables -L OUTPUT --line-numbers
iptable -D OUTPUT 1
```

2. 查看切换日志

```
ansible -i ./zerodb-ansible/hosts keeper_daily -m "shell" -a "tail -200 /var/log/zlog/ds_switch.log.2018082900" --extra-vars "ansible_user=root ansible_password=sdfsdf"

ansible -i ./zerodb-ansible/hosts 10.1.21.82 -m "shell" -a "tail -200 /var/log/zlog/ds_switch.log.2018082900" --extra-vars "ansible_user=root ansible_password=sdfsdf"
```

### 2. 无法连接多台mysql,注销自己 

**预期结果**: proxy 因为因为跟多台mysql无法连接,向keeper注销自己、mysql引擎关闭、checkhealth server关闭、 etcd中proxy被删除

1. 设置proxy 10.1.21.82 iptables 无法访问多个mysql 3306

```
iptables -A OUTPUT  -p tcp  -d 10.1.22.3  --dport 3306 -j REJECT
iptables -A OUTPUT  -p tcp  -d 10.1.22.1  --dport 3306 -j REJECT
```

2. 查看proxy 状态日志和切换日志

```
ansible -i ./zerodb-ansible/hosts 10.1.21.82 -m "shell" -a "cat /var/log/zlog/proxy.log.2018083000" --extra-vars "ansible_user=root ansible_password=sdfsdf" |grep "offline" 

ansible -i ./zerodb-ansible/hosts proxy_daily -m "shell" -a "tail -200 /var/log/zlog/ds_switch.log.2018083000" --extra-vars "ansible_user=root ansible_password=sdfsdf" 
```

3. 查看etcd中proxy的状态、mysql引擎状态、checkhealth server状态

```
etcdctl --endpoints=10.1.24.173:2379 get --prefix zerodb/proxy/status/cluster-func-test/
telnet 10.1.21.82 9696 9002
```

### 4.proxy网络恢复,重新注册 

**预期结果**: proxy因为跟所有mysql重新连接,重现恢复自己、mysql 引擎启动、etcd中成功注册.

1. 删除iptables,恢复proxy 10.1.21.82网络

```
iptables -L OUTPUT --line-numbers
iptables -D 1
```

2. 查看proxy 状态日志

```
ansible -i ./zerodb-ansible/hosts 10.1.21.82 -m "shell" -a "cat /var/log/zlog1/proxy.log.2018082900" --extra-vars "ansible_user=root ansible_password=sdfsdf" |grep "recover" 
```

3. 查看etcd中proxy的状态、mysql引擎状态、checkhealth server状态

```
etcdctl --endpoints=10.1.24.173:2379 get --prefix zerodb/proxy/status/cluster-func-test/
telnet 10.1.21.82 9696 9002
```

## 5. 所有proxy都无法连接到mysql

1. 设置iptables 所有ip都无法访问10.1.22.3 和10.1.22.1 的3306端口

```
iptables -I INPUT -p tcp -s 10.1.21.79 --dport 3306 -j DROP
iptables -I INPUT -p tcp -s 10.1.21.82 --dport 3306 -j DROP
iptables -I INPUT -p tcp -s 10.1.21.242 --dport 3306 -j DROP
iptables -L INPUT --line-numbers
iptables -D INPUT 1
```

2. 查看proxy日志

```
ansible -i ./zerodb-ansible/hosts proxy_daily -m "shell" -a "cat /var/log/zlog/proxy.log.2018083000" --extra-vars "ansible_user=root ansible_password=sdfsdf" |grep "offline" 
```

3. 查看etcd中的proxy状态

```
etcdctl --endpoints=10.1.24.173:2379 get --prefix zerodb/proxy/status/cluster-func-test/
```

4. 删除iptables,恢复所有的proxy

```
iptables -D INPUT 1

ansible -i ./zerodb-ansible/hosts 10.1.21.82 -m "shell" -a "cat /var/log/zlog/proxy.log.2018082900" --extra-vars "ansible_user=root ansible_password=sdfsdf" |grep "recover" 
```



## 6. keeper 挂掉 

1. kill  zeorkeeper

```
ansible -i ./zerodb-ansible/hosts keeper_daily -m "shell" -a " ps -ef |grep zero " --extra-vars "ansible_user=root ansible_password=sdfsdf"

ansible -i ./zerodb-ansible/hosts keeper_daily -m "shell" -a "kill -9 17433 " --extra-vars "ansible_user=root ansible_password=sdfsdf"
```

2. 查看proxy状态日志 进程 

```
ansible -i ./zerodb-ansible/hosts proxy_daily -m "shell" -a "ls /var/log/zlog" --extra-vars "ansible_user=root ansible_password=sdfsdf"

ansible -i ./zerodb-ansible/hosts proxy_daily -m "shell" -a "tail -100 /var/log/zlog/heartbeat.log.2018082900" --extra-vars "ansible_user=root ansible_password=sdfsdf"

ansible -i ./zerodb-ansible/hosts proxy_daily -m "shell" -a " ps -ef |grep zero " --extra-vars "ansible_user=root ansible_password=sdfsdf"
```



## 7. agent 挂掉

1. kill mysql 10.1.21.81  agent 

  ```
  ansible -i ./zerodb-ansible/hosts 10.1.21.81 -m "shell" -a " ps -ef |grep zero " --extra-vars "ansible_user=root ansible_password=sdfsdf"
  ```

```
ansible -i ./zerodb-ansible/hosts 10.1.21.81 -m "shell" -a "kill -9 pid " --extra-vars "ansible_user=root ansible_password=sdfsdf"
```

2. 设置iptables mysql从10.1.22.3 切换到10.1.21.81

```
iptables -I INPUT -p tcp -s 10.1.21.79 --dport 3306 -j DROP
iptables -I INPUT -p tcp -s 10.1.21.82 --dport 3306 -j DROP
```

3. 查看切换日志

```\
ansible -i ./zerodb-ansible/hosts keeper_daily -m "shell" -a "tail -200 /var/log/zlog/ds_switch.log.2018082900" --extra-vars "ansible_user=root ansible_password=sdfsdf"

ansible -i ./zerodb-ansible/hosts proxy_daily -m "shell" -a "tail -200 /var/log/zlog/ds_switch.log.2018082900" --extra-vars "ansible_user=root ansible_password=sdfsdf"

```

4. 跳过agent check 和binlogcheck。

## 8.请求锁

**预期结果**: 只有一台proxy获得锁.

```
ansible -i ./zerodb-ansible/hosts keeper_daily -m "shell" -a "cat  /var/log/zlog/ds_switch.log.2018083000" --extra-vars "ansible_user=root ansible_password=sdfsdf" |grep "lock"

ansible -i ./zerodb-ansible/hosts proxy_daily -m "shell" -a "cat  /var/log/zlog/ds_switch.log.2018083000" --extra-vars "ansible_user=root ansible_password=sdfsdf" |grep "lock"
```

## 9. 切换频率

**预期结果**: 30秒内一个hostgroup只被允许切换一次

```
ansible -i ./zerodb-ansible/hosts keeper_daily -m "shell" -a "cat  /var/log/zlog/ds_switch.log.2018083000" --extra-vars "ansible_user=root ansible_password=sdfsdf" |grep "lock"


```

## 10. binlog 同步

**预期结果**: binlog 不同步无法切换

## 11. Proxy 状态 



### 

### 

