package com.dfire.zerodb.simeple;


import com.dfire.zerodb.common.DbPoolConnection;
import com.alibaba.druid.pool.DruidPooledConnection;

import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;

/**
 * User: nanxing
 * Date: 11/2/16
 */
public class SimpleDruidTest {

    @org.junit.Test
    public void fetchDatasource() throws SQLException {
        DbPoolConnection dbp = DbPoolConnection.getInstance();
        DruidPooledConnection con = dbp.getConnection();

        String sql = "select * from user";

        PreparedStatement ps = con.prepareStatement(sql);
        ResultSet resultSet = ps.executeQuery();
        while (resultSet.next()) {
//            String string = resultSet.getString(4);
//            System.out.println(string);
            System.out.println("id:" + resultSet.getString(1) + ",member_id:" + resultSet.getString(2));
        }
        ps.close();
        con.close();
        dbp = null;
    }
}
