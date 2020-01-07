## 文件

### proxy:

#### heartbeat.log
> ->keeper的心跳

#### ds_switch.log
> proxy数据源切换

#### slow_query.log
> 慢SQL日志

#### client_conn.log
> app->proxy连接的建立与断开

#### api.log
> proxy端的web-api调用

#### proxy.log
> 启动与逻辑

#### watermark.log
> 水位监控与控制

#### dangerous_sql.log
> 多路由SQL的限制与放行

#### db_pool.log
> proxy->MySQL数据库连接池相关

### keeper:

#### ds_switch.log
> 数据源切换决策过程

#### proxy_status.log
> proxy状态转换过程

#### api_handle.log
> web api调用

#### proxy_hb.log
> proxy->的心跳

#### agent_hb.log
> agent->的心跳

#### keeper.log
> 启动与逻辑

### agent:

heartbeat.log
> ->keeper的心跳

agent.log
> 启动与逻辑

#### api.log
> agent端的web-api调用

## 格式:

```text
B: for other log
2018/05/31 13:31:44.995 2pc.go:582: [info] 2PC clean up done, tid: 400489107892994203

C: for hb
2018/05/31 13:31:44.995 [info] 2PC clean up done, tid: 400489107892994203
```
