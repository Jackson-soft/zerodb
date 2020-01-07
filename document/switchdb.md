#  ZeroKeeper 数据源切换

   ​

   ​

   ## zerokeeper 数据源切换接口实现

   ### 参数

   * Hostgroup    // hostgroup Name
   * From.            // 需要切换的mysql host索引 (当前连接的) mysql host 索引
   * To.                // 切换到的mysql host 索引
   * ProxyIP        // 发起切换的ZeroProxy IP.
   * ClusterName  //ZeroProxy 的集群名字

   数据源切换相关的配置

   * NeedVote   //切换时是否需要投票,true 允许、false 不允许.
   * VoteApproveRatio  //投票比列, 超过这个比例即认为集群中所有zeroproxy允许切换
   * NeedLoadCheck  //是否需要检查mysql load,true 允许、false 不允许
   * SafeLoad  //安全的load值,超过改值不允许切换
   * NeedBinlogCheck  //是否需要检查binlog
   * SafeBinlogDelay  //安全的binlog 延迟数值.

   ### 切换过程

   当zerokeeper 收到来自zeroproxy的数据源切换请求时,会根据一些条件判断来决定是否允许此次数据源切换.

   1. 根据ClusterName参数判断该集群中所有zeroproxy的状态,如果有任何一台zeroproxy处于非up状态,拒绝切换返回错误.

   2. 判断切换的频率,对于一台hostgroup只允许30秒内切换一次.

   3. 检查mysql负载

   4. 获取分布式锁,任何proxy的切换请求都需要获得锁才能被处理.无法获取锁直接返回.锁根据hostgroup name建立在etcd中

   5. 停止所有proxy的写操作.

   6. 检查from 和 to 的 binlog位置差距.如果binlog 位置差距大于配置的安全延迟值等待3秒钟binlog同步,如果仍然大于安全延迟值,返回.

   7. 获取所有zeroproxy的投票, 超过投票比例即认为投票通过,任何proxy投票返回错误,直接返回.

   8. 发送数据源切换给所有zeroproxy,任何zeroproxy返回错误直接返回

      ​

