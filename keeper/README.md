zeroKeeper是zeroProxy的集群管理者与协调者，功能比较复杂

zeroKeeper<->zeroProxy:
协调proxy集群数据源切换、协调配置下发

zeroKeeper<->MySQL:
数据库连接池、MySQL心跳存活检测

zeroKeeper<->zeroAgent:
拉取agent管理的MySQL的负载