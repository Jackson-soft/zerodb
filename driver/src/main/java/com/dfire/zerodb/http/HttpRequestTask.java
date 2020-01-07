package com.dfire.zerodb.http;

import com.dfire.zerodb.common.ProxyNode;

import java.util.List;
import java.util.concurrent.CompletableFuture;

public class HttpRequestTask implements Runnable {
    private CompletableFuture<GetProxiesResponse> future;
    private String urlPath;
    private String clusterName;

    public HttpRequestTask(CompletableFuture<GetProxiesResponse> future, String urlPath, String clusterName) {
        this.future = future;
        this.urlPath = urlPath;
        this.clusterName = clusterName;
    }

    @Override
    public void run() {
        try {
            List<ProxyNode> proxyNodes = HttpRequester.doRequest(urlPath, clusterName);
            future.complete(new GetProxiesResponse(true, proxyNodes));
        } catch (Exception e) {
            future.completeExceptionally(e);
        }
    }
}
