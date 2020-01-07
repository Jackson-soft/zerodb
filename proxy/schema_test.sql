
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

delete from order_ins where entity_id = 123125;
INSERT INTO order_ins (id, entity_id, name) VALUES (6, 123125, 'hello world🌀order_ins');
select * from order_ins where entity_id = 123125;

delete from ins_detail where entity_id = 123125;
INSERT INTO ins_detail (id, entity_id, detail) VALUES (6, 123125, 'hello world👿');
select * from ins_detail where entity_id = 123125;

select id.id, id.detail from ins_detail id inner join order_ins oi on oi.entity_id = id.entity_id where oi.entity_id = 123125;
select id.id, id.detail from ins_detail id, order_ins oi where oi.entity_id = id.entity_id and oi.entity_id = 123125;


## transaction test
use zerodb;
set autocommit = 1;
delete from order_ins where entity_id = 123125;
delete from ins_detail where entity_id = 123125;
set autocommit = 0;
INSERT INTO order_ins (id, entity_id, name) VALUES (6, 123125, 'hello world🌀order_ins');
INSERT INTO order_ins (id, entity_id, name) VALUES (6, 123126, 'hello world🌀order_ins');
INSERT INTO ins_detail (id, entity_id, detail) VALUES (6, 123125, 'hello world👿');
commit;
select * from order_ins where entity_id = 123125;
select * from ins_detail where entity_id = 123125;

## transaction test
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

## transaction test in nonsharding node
use zerodb;
set autocommit = 1;
delete from order_ins where entity_id = 123125;
delete from ins_detail where entity_id = 123125;
set autocommit = 0;
INSERT INTO order_ins (id, entity_id, name) VALUES (6, 123125, 'hello world🌀order_ins');
select * from order_ins where entity_id = 123;
INSERT INTO ins_detail (id, entity_id, detail) VALUES (6, 123125, 'hello world👿');
rollback;

## transaction test in custody schema
use zerodb;
rollback;

## multinodes insert
## 不支持
INSERT INTO order_ins (id, entity_id, name) VALUES (6, 123125, 'hello world🌀order_ins1'), (7, 123126, 'hello world🌀order_ins2');


CREATE TABLE `card` (
    `id`        INT(11) NOT NULL,
    `entity_id` BIGINT(20)  DEFAULT NULL,
    `name`      VARCHAR(45) DEFAULT NULL,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO card (id, entity_id, name) VALUES (6, 123125, 'hello');