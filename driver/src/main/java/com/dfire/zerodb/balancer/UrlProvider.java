package com.dfire.zerodb.balancer;


import com.dfire.zerodb.common.ProxyNode;
import com.dfire.zerodb.http.GetProxiesResponse;
import com.dfire.zerodb.http.HttpClient;

import java.sql.SQLException;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

public class UrlProvider {
    private static final String URL_ZERODB_PREFIX = "jdbc:zerodb://";
    private static final String URL_MYSQL_PREFIX = "jdbc:mysql://";
    private static final int URL_ZERODB_PREFIX_LENGTH = URL_ZERODB_PREFIX.length();
    private static final int REQUEST_TIMEOUT = 3000;
    private static final String[] EMPTY_STRING_ARRAY = new String[0];
    private static final Random RANDOM = new Random();

    /**
     * 只处理非空并且以jdbc:zerodb:/开头的url
     */
    public static final String getUrl(String url, Properties info) throws SQLException {
        if (url != null && url.regionMatches(true, 0, URL_ZERODB_PREFIX, 0, URL_ZERODB_PREFIX_LENGTH)) {
            try {
                ConnectInfo ci = parseUrl(url, info);
                final CompletableFuture<GetProxiesResponse> future = HttpClient.request(ci.host, ci.clusterName);
                final GetProxiesResponse response = future.get(REQUEST_TIMEOUT, TimeUnit.MILLISECONDS);
                return selectUrl(response.getProxyNodes(), url, ci);
            } catch (Throwable e) {
                throw new SQLException(e);
            }
        } else {
            return url;
        }
    }

    /**
     * 取得MySQL格式的URL
     */
    public static final String getMySQLUrl(String url) {
        if (url != null && url.regionMatches(true, 0, URL_ZERODB_PREFIX, 0, URL_ZERODB_PREFIX_LENGTH)) {
            StringBuilder sb = new StringBuilder();
            sb.append(URL_MYSQL_PREFIX).append(url.substring(URL_ZERODB_PREFIX_LENGTH));
            url = sb.toString();
        }
        return url;
    }

    private static ConnectInfo parseUrl(String url, Properties info) {
        ConnectInfo ci = new ConnectInfo();
        ci.user = info.getProperty("user");
        ci.password = info.getProperty("password");
        ci.database = info.getProperty("DBNAME");
        ci.clusterName = info.getProperty("clusterName");

        // 解析参数
        url = parseParameter(url, ci);

        // 解析host/port/database
        url = url.substring(URL_ZERODB_PREFIX_LENGTH);
        String hostStuff = null;
        int index2 = url.indexOf('/');
        if (index2 != -1) {
            hostStuff = url.substring(0, index2);
            if (index2 + 1 < url.length()) {
                ci.database = url.substring(index2 + 1);
            }
        } else {
            hostStuff = url;
        }
        // 如果有多个主机只取第一个
        int index3 = hostStuff.indexOf(',');
        if (index3 != -1) {
            hostStuff = hostStuff.substring(0, index3);
        }
        int index4 = hostStuff.indexOf(':');
        if (index4 != -1) {
            ci.host = hostStuff.substring(0, index4).trim();
        } else {
            ci.host = hostStuff.trim();
        }

        return ci;
    }

    private static String parseParameter(String url, ConnectInfo ci) {
        int index = url.indexOf('?');
        if (index != -1) {
            String paramString = url.substring(index + 1);
            url = url.substring(0, index);
            ci.paramString = paramString;

            // ?后面的参数可以覆盖properties配置
            String[] params = split(paramString, '&');
            for (String param : params) {
                int indexOfEquals = 0;
                if (param != null && (indexOfEquals = param.indexOf('=')) != -1) {
                    String key = param.substring(0, indexOfEquals);
                    if ("user".equals(key)) {
                        if (indexOfEquals + 1 < param.length()) {
                            ci.user = param.substring(indexOfEquals + 1);
                        }
                    } else if ("password".equals(key)) {
                        if (indexOfEquals + 1 < param.length()) {
                            ci.password = param.substring(indexOfEquals + 1);
                        }
                    } else if ("DBNAME".equals(key)) {
                        if (indexOfEquals + 1 < param.length()) {
                            ci.database = param.substring(indexOfEquals + 1);
                        }
                    } else if ("clusterName".equals(key)) {
                        if (indexOfEquals + 1 < param.length()) {
                            ci.clusterName = param.substring(indexOfEquals + 1);
                        }
                    }
                }
            }
        }
        return url;
    }

    /**
     * 结合权重和随机策略
     */
    private static String selectUrl(List<ProxyNode> list, String originUrl, ConnectInfo info) {
        int total = 0;
        for (ProxyNode node : list) {
            total += node.getWeight();
        }
        // 如果总权重小于等于零则使用原来的url
        if (total <= 0) {
            return originUrl;
        }
        int rnd = 1 + RANDOM.nextInt(total);
        ProxyNode selected = null;
        for (ProxyNode node : list) {
            if ((rnd -= node.getWeight()) <= 0) {
                selected = node;
                break;
            }
        }
        if (selected == null) {
            int i = RANDOM.nextInt(list.size());
            selected = list.get(i);
        }

        // 无法取得新的proxy，则返回原来的url。
        if (selected == null) {
            return originUrl;
        }

        if ("".equals(selected.getHost())) {
            throw new RuntimeException("Host of proxy node is impty.");
        }
        if (selected.getPort() == 0) {
            throw new RuntimeException("Port of proxy node is zero.");
        }

        // 生成新的URL
        StringBuilder url = new StringBuilder();
        url.append(URL_MYSQL_PREFIX).append(selected.getHost()).append(':').append(selected.getPort());
        if (info.database != null) {
            url.append('/').append(info.database);
        }
        if (info.paramString != null) {
            url.append('?').append(info.paramString);
        }
        return url.toString();
    }

    private static String[] split(String src, char separatorChar) {
        if (src == null) {
            return null;
        }
        int length = src.length();
        if (length == 0) {
            return EMPTY_STRING_ARRAY;
        }
        List<String> list = new LinkedList<String>();
        int i = 0;
        int start = 0;
        boolean match = false;
        while (i < length) {
            if (src.charAt(i) == separatorChar) {
                if (match) {
                    list.add(src.substring(start, i));
                    match = false;
                }
                start = ++i;
                continue;
            }
            match = true;
            i++;
        }
        if (match) {
            list.add(src.substring(start, i));
        }
        return list.toArray(new String[list.size()]);
    }

    private static class ConnectInfo {
        private String host;
        private String user;
        private String password;
        private String database;
        private String paramString;
        private String clusterName;
    }

}
