#!/usr/local/bin/python3
# -*- coding: utf-8 -*-

import pymysql

# 打开数据库连接
# db = pymysql.connect("common101.my.2dfire-daily.com", "guard", "guard@552208", "guard")
db = pymysql.connect(
    host="10.12.34.23",
    port=5003,
    user="zerodb",
    password="zerodb",
    db="guard",
    charset="utf8",
)

# 使用 cursor() 方法创建一个游标对象 cursor
cursor = db.cursor()

# 使用 execute()  方法执行 SQL 查询
cursor.execute("SELECT VERSION()")

# 使用 fetchone() 方法获取单条数据.
data = cursor.fetchone()

print("Database version : %s " % data)

# 关闭数据库连接
db.close()
