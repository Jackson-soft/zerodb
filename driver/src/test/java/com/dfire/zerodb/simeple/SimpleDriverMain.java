/*
 * Copyright 1999-2012 Alibaba Group.
 *  
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *  
 *      http://www.apache.org/licenses/LICENSE-2.0
 *  
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package com.dfire.zerodb.simeple;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.ResultSet;
import java.sql.Statement;
import java.util.Properties;

/**
 * @author xianmao.hexm 2012-4-27
 */
public class SimpleDriverMain {

    public static void main(String[] args) throws Exception {
//        Class.forName("com.dfire.zerodb.driver.Driver");
        //String url = "jdbc:zerodb://zerodb.app.2dfire-daily.com/missile";
        String url = "jdbc:zerodb://zerodb.app.2dfire-daily.com/missile";
        Properties info = new Properties();
        info.setProperty("user", "zerodb");
        info.setProperty("password", "zerodb");
        info.setProperty("clusterName", "daily");
        Connection con = null;
        try {
            con = DriverManager.getConnection(url, info);
            Statement stmt = con.createStatement();
//            String query = "INSERT INTO `user` (`user_id`, `name`, `birthday`, `gender`, `region_id`, `intro`) VALUES ('12122121213213','0005748756961e470156b2ef80786684\u0001香烤草鱼\u0001\uD83D\uDC1F', '213', '123', '123', '123')INSERT INTO `user` (`user_id`, `name`, `birthday`, `gender`, `region_id`, `intro`) VALUES ('12122121213213','0005748756961e470156b2ef80786684\u0001香烤草鱼\u0001\uD83D\uDC1F', '213', '123', '123', '123');";
//            String query = "select * from instancedetail limit 1;";
            String query = "select * from user";
            ResultSet rs = stmt.executeQuery(query);
            while (rs.next()) {
                System.out.println("id:" + rs.getString(1) + ",member_id:" + rs.getString(2));
            }
            rs.close();
            stmt.close();
        } catch (Exception e) {
            System.out.println(e);
        } finally {
            if (con != null) {
                con.close();
            }
        }
    }

}
