package com.dfire.zerodb.http;


import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Executor;
import java.util.concurrent.Executors;

public class HttpClient {
    private static String threadNamePrefix = "HTTP-Client";

    private static Executor threadPool = Executors.newFixedThreadPool(1, r -> {
        Thread thread = new Thread(r);
        thread.setDaemon(true);
        thread.setName(threadNamePrefix + thread.getId());
        return thread;
    });

    public static CompletableFuture<GetProxiesResponse> request(String url, String clusterName) {
        CompletableFuture future = new CompletableFuture();
        final HttpRequestTask requestTask = new HttpRequestTask(future, url, clusterName);
        threadPool.execute(requestTask);
        return future;
    }
}
