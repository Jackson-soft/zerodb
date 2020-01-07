南星（半小时）

## 分库分表演示

#### 分库(0+128)--CREATE(DATABASE, TABLE), SELECT(TABLE)
```
show databases;

use zerodb;
drop database zerodb;

create database zerodb;

CREATE TABLE `order_ins` (
    `id`        INT(11) NOT NULL,
    `entity_id` BIGINT(20)  DEFAULT NULL,
    `name`      VARCHAR(45) DEFAULT NULL,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `ins_detail` (
    `id`        INT(11) NOT NULL,
    `entity_id` BIGINT(20)  DEFAULT NULL,
    `detail`      VARCHAR(45) DEFAULT NULL,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `order_ins_detail` (
    `id`        INT(11) NOT NULL,
    `entity_id` BIGINT(20)  DEFAULT NULL,
    `detail`      VARCHAR(45) DEFAULT NULL,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `non_sharding_table` (
    `id`        INT(11) NOT NULL,
    `entity_id` BIGINT(20)  DEFAULT NULL,
    `detail`      VARCHAR(45) DEFAULT NULL,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

route create database zerodb;

route select * from order_ins;
route select * from ins_detail;
route select * from order_ins_detail;

delete from order_ins where entity_id = 123125;
INSERT INTO order_ins (id, entity_id, name) VALUES (6, 123125, 'hello world🌀order_ins');
select * from order_ins where entity_id = 123125;

delete from ins_detail where entity_id = 123125;
INSERT INTO ins_detail (id, entity_id, detail) VALUES (6, 123125, 'hello world👿');
select * from ins_detail where entity_id = 123125;

INSERT INTO ins_detail (id, entity_id, detail) VALUES (1, 11111111, 'hello world👿'),(2, 222222, 'hello world👿');

select id.id, id.detail from ins_detail id inner join order_ins oi on oi.entity_id = id.entity_id where oi.entity_id = 123125;
select id.id, id.detail from ins_detail id, order_ins oi where oi.entity_id = id.entity_id and oi.entity_id = 123125;

```

#### 分库(1+128)--CREATE(DATABASE, TABLE), SELECT(TABLE)
```
show databases;

use instance;
drop database instance;

create database instance;

CREATE TABLE `order_ins` (
    `id`        INT(11) NOT NULL,
    `entity_id` BIGINT(20)  DEFAULT NULL,
    `name`      VARCHAR(45) DEFAULT NULL,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `ins_detail` (
    `id`        INT(11) NOT NULL,
    `entity_id` BIGINT(20)  DEFAULT NULL,
    `detail`      VARCHAR(45) DEFAULT NULL,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `order_ins_detail` (
    `id`        INT(11) NOT NULL,
    `entity_id` BIGINT(20)  DEFAULT NULL,
    `detail`      VARCHAR(45) DEFAULT NULL,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `non_sharding_table` (
    `id`        INT(11) NOT NULL,
    `entity_id` BIGINT(20)  DEFAULT NULL,
    `detail`      VARCHAR(45) DEFAULT NULL,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

route create database zerodb;

route select * from order_ins;
route select * from ins_detail;
route select * from order_ins_detail;
route select * from non_sharding_table;

delete from order_ins where entity_id = 123125;
INSERT INTO order_ins (id, entity_id, name) VALUES (6, 123125, 'hello world🌀order_ins');
select * from order_ins where entity_id = 123125;

delete from ins_detail where entity_id = 123125;
INSERT INTO ins_detail (id, entity_id, detail) VALUES (6, 123125, 'hello world👿');
select * from ins_detail where entity_id = 123125;

```

#### 分库+分表(1+128*8)--CREATE(DATABASE, TABLE), SELECT(TABLE)
```
show databases;

use many;
drop database many;

create database many;

CREATE TABLE `order_ins` (
    `id`        INT(11) NOT NULL,
    `entity_id` BIGINT(20)  DEFAULT NULL,
    `name`      VARCHAR(45) DEFAULT NULL,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `ins_detail` (
    `id`        INT(11) NOT NULL,
    `entity_id` BIGINT(20)  DEFAULT NULL,
    `detail`      VARCHAR(45) DEFAULT NULL,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `order_ins_detail` (
    `id`        INT(11) NOT NULL,
    `entity_id` BIGINT(20)  DEFAULT NULL,
    `detail`      VARCHAR(45) DEFAULT NULL,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `non_sharding_table` (
    `id`        INT(11) NOT NULL,
    `entity_id` BIGINT(20)  DEFAULT NULL,
    `detail`      VARCHAR(45) DEFAULT NULL,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

delete from order_ins where entity_id = 123125;
INSERT INTO order_ins (id, entity_id, name) VALUES (6, 123125, 'hello world🌀order_ins');
select * from order_ins where entity_id = 123125;

delete from ins_detail where entity_id = 123125;
INSERT INTO ins_detail (id, entity_id, detail) VALUES (6, 123125, 'hello world👿');
select * from ins_detail where entity_id = 123125;


select id.id, id.detail from ins_detail id inner join order_ins oi on oi.entity_id = id.entity_id where oi.entity_id = 123125;
select id.id, id.detail from ins_detail id, order_ins oi where oi.entity_id = id.entity_id and oi.entity_id = 123125;
```

#### 不分库不分表——托管(1*1)--CREATE(DATABASE, TABLE), SELECT(TABLE)
```
use custody;
drop database custody;

create database custody;

CREATE TABLE `order_ins` (
    `id`        INT(11) NOT NULL,
    `entity_id` BIGINT(20)  DEFAULT NULL,
    `name`      VARCHAR(45) DEFAULT NULL,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `ins_detail` (
    `id`        INT(11) NOT NULL,
    `entity_id` BIGINT(20)  DEFAULT NULL,
    `detail`      VARCHAR(45) DEFAULT NULL,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

route select * from order_ins;
route select * from ins_detail;

```
#### 读写分离
```
"custody": true
```

```
use custody;
INSERT INTO order_ins (id, entity_id, name) VALUES (6, 123125, 'hello world🌀order_ins'); // 有bug

delete from order_ins where entity_id = 123125;
select * from order_ins where entity_id = 123125;
```

#### 非跨库事务（影响性能）
```

use zerodb;
set autocommit = 1;
delete from order_ins where entity_id = 123125;
delete from ins_detail where entity_id = 123125;
set autocommit = 0;
INSERT INTO order_ins (id, entity_id, name) VALUES (6, 123125, 'hello world🌀order_ins');

INSERT INTO ins_detail (id, entity_id, detail) VALUES (6, 123125, 'hello world👿');
select * from order_ins where entity_id = 123125;
select * from ins_detail where entity_id = 123125;
commit;
select * from order_ins where entity_id = 123125;
select * from ins_detail where entity_id = 123125;

use zerodb;
set autocommit = 1;
delete from order_ins where entity_id = 123125;
delete from ins_detail where entity_id = 123125;
set autocommit = 0;
INSERT INTO order_ins (id, entity_id, name) VALUES (6, 123125, 'hello world🌀order_ins');
select * from order_ins where entity_id = 123125;
INSERT INTO ins_detail (id, entity_id, detail) VALUES (6, 123125, 'hello world👿');
select * from ins_detail where entity_id = 123125;
rollback;
select * from order_ins where entity_id = 123125;
select * from ins_detail where entity_id = 123125;

// 扩分片
use zerodb;
set autocommit = 1;
delete from order_ins where entity_id = 123125;
delete from ins_detail where entity_id = 666125;
set autocommit = 0;
INSERT INTO order_ins (id, entity_id, name) VALUES (6, 123125, 'hello world🌀order_ins');
INSERT INTO ins_detail (id, entity_id, detail) VALUES (6, 666125, 'hello world👿');
select * from order_ins where entity_id = 123125;
select * from ins_detail where entity_id = 666125;
rollback;
select * from order_ins where entity_id = 123125;
select * from ins_detail where entity_id = 666125;
```

#### Route语句
```
use many
route create database zerodb;
route drop database zerodb;
route select * from order_ins;
route select * from order_ins where entity_id in (1, 33333, 444444);
route select * from order_ins where entity_id = 123;
route select id.id, id.detail from ins_detail id, order_ins oi where oi.entity_id = id.entity_id and oi.entity_id = 123125\G
route delete from order_ins; # 没有事务
route delete from order_ins where entity_id = 123;
route update order_ins set name = 'aaaa'; # 没有事务
route update order_ins set name = 'aaaa' where entity_id = 12223;
route INSERT INTO order_ins (id, entity_id, name) VALUES (6, 123, 'hello world🌀order_ins');
INSERT INTO ins_detail (id, entity_id, detail) VALUES (1, 11111111, 'hello world👿'),(2, 222222, 'hello world👿');
INSERT INTO ins_detail (id, entity_id, detail) VALUES (1, 1, 'hello world👿'),(2, 2, 'hello world👿'),(3, 222222, 'hello world👿');
route drop table order_ins;
```

#### 限制MultiRouteSQL
```
set MULTI_ROUTE_PERMIT = 0;
select * from order_ins;
set MULTI_ROUTE_PERMIT = 1;
select * from order_ins;
```

## 性能指标
![image.png](https://upload-images.jianshu.io/upload_images/716353-ab539e0b6d2f4339.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

![image.png](https://upload-images.jianshu.io/upload_images/716353-8063cb58a9b6093d.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

![image.png](https://upload-images.jianshu.io/upload_images/716353-ba90d4411a1fe7a5.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

![image.png](https://upload-images.jianshu.io/upload_images/716353-300e7a438ad3df91.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

## 局限性
    1. 弱一致性
    2. 只支持非扩库单Join语句 select a.x, b.y from a join b on a.id = b.id
    3. 单路由键(int, string)
## 优点
    1. 智能水位控制
    2. 智能数据源切换
    3. 集群配置管理
    4. 读写分离
    5. 集群Proxy
    6. 智能恢复机制
    7. 分库+分表
    8. QPS/TPS统计

泡面（半小时）
1. 分布式调度
    数据源切换--投票
    数据源切换--binlog
    数据源切换--load
2. StopService

扶苏（半小时）
1. API演示
### snapshot创建(每一步都要证明是成功了，还是失败了)
    1. Proxy(3 live)
    2. Proxy(2 live + 1 dead)
    3. Proxy(3 dead)

### 推送snapshot
    1. Proxy(3 live)
    2. Proxy(2 live + 1 dead)
    3. Proxy(3 dead)


### snapshot分发
2. Proxy状态转移
    proxy-status-transfer.md里面的过程走一遍