package com.dfire.zerodb.common;

public class ProxyNode {

    private String host;
    private int port;
    private int weight;

    public ProxyNode(String host, int port, int weight) {
        this.host = host;
        this.port = port;
        this.weight = weight;
    }

    public String getHost() {
        return host;
    }

    public int getPort() {
        return port;
    }

    public int getWeight() {
        return weight;
    }
}
