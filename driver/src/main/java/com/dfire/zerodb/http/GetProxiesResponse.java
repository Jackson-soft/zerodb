package com.dfire.zerodb.http;

import com.dfire.zerodb.common.ProxyNode;

import java.util.List;

public class GetProxiesResponse {
    private boolean success;
    private List<ProxyNode> proxyNodes;

    public GetProxiesResponse(boolean success, List<ProxyNode> proxyNodes) {
        this.success = success;
        this.proxyNodes = proxyNodes;
    }

    public boolean isSuccess() {
        return success;
    }

    public List<ProxyNode> getProxyNodes() {
        return proxyNodes;
    }
}
