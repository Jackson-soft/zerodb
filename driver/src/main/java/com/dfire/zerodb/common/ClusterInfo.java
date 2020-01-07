package com.dfire.zerodb.common;

import java.util.List;

public class ClusterInfo {
    private int code;
    private int errorCode;
    private List<ProxyNode> data;

    public ClusterInfo(int code, int errorCode, List<ProxyNode> data) {
        this.code = code;
        this.errorCode = errorCode;
        this.data = data;
    }

    public int getCode() {
        return code;
    }

    public int getErrorCode() {
        return errorCode;
    }

    public List<ProxyNode> getData() {
        return data;
    }
}
