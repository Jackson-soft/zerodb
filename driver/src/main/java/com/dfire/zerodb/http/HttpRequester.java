package com.dfire.zerodb.http;

import com.dfire.zerodb.common.ClusterInfo;
import com.dfire.zerodb.common.ProxyNode;
import com.alibaba.fastjson.JSONObject;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.List;

public class HttpRequester {
    static List<ProxyNode> doRequest(String urlPath, String clusterName) {
        HttpURLConnection connection = null;
        String u = "https://" + urlPath + "/proxy_list?cluster_name=" + clusterName;

        try {
            URL url = new URL(u);
            connection = (HttpURLConnection) url.openConnection();

            connection.setDoOutput(true); // 设置该连接是可以输出的
            connection.setRequestMethod("GET"); // 设置请求方式
            connection.setRequestProperty("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8");

            BufferedReader br = new BufferedReader(new InputStreamReader(connection.getInputStream(), "utf-8"));
            String line = null;
            StringBuilder result = new StringBuilder();
            while ((line = br.readLine()) != null) { // 读取数据
                result.append(line + "\n");
            }
            connection.disconnect();

            ClusterInfo clusterInfo = JSONObject.parseObject(result.toString(), ClusterInfo.class);

            return clusterInfo.getData();
        } catch (Exception e) {
            throw new RuntimeException("Request Cluster List Error: [" + u + "]" + e);
        } finally {
            if (connection != null) {
                connection.disconnect();
            }
        }
    }
}
