package com.dfire.zerodb.driver;

import com.dfire.zerodb.balancer.UrlProvider;
import com.mysql.jdbc.NonRegisteringDriver;

import java.sql.Connection;
import java.sql.DriverPropertyInfo;
import java.sql.SQLException;
import java.sql.SQLFeatureNotSupportedException;
import java.util.Properties;
import java.util.logging.Logger;

/**
 * 在使用集群时提供负载均衡的功能，其他情况和MySQLDriver一样。
 *
 * <pre>
 * 使用方法：
 *   Class.forName("com.dfire.zerodb.driver.Driver");
 *   String url = "jdbc:zerodb://host/dbname?user=xxx&password=xxx";
 *   ...
 * </pre>
 *
 * @author nanxing
 */
public class Driver extends NonRegisteringDriver implements java.sql.Driver {

    public static final String VERSION = "1.0.0";

    /**
     * Register ourselves with the DriverManager
     */
    static {
        try {
            java.sql.DriverManager.registerDriver(new Driver());
        } catch (SQLException E) {
            throw new RuntimeException("Can't register driver!");
        }
    }

    /**
     * Construct a new driver and register it with DriverManager
     *
     * @throws SQLException if a database error occurs.
     */
    public Driver() throws SQLException {
        // Required for Class.forName().newInstance()
    }

    @Override
    public Connection connect(String url, Properties info) throws SQLException {
        String s = null;
        try {
            s = UrlProvider.getUrl(url, info);
            return super.connect(s, info);
        } catch (SQLException e) {
            throw new SQLException("Real [url: " + s + "] connection failed. Reason: " + e.getMessage());
        }
    }

    @Override
    public boolean acceptsURL(String url) throws SQLException {
        return super.acceptsURL(UrlProvider.getMySQLUrl(url));
    }

    @Override
    public DriverPropertyInfo[] getPropertyInfo(String url, Properties info) throws SQLException {
        return super.getPropertyInfo(UrlProvider.getMySQLUrl(url), info);
    }

    //@Override
    public Logger getParentLogger() throws SQLFeatureNotSupportedException {
        throw new SQLFeatureNotSupportedException("Not supported yet.");
    }

}
