package com.dfire.zerodb.common;

import java.util.concurrent.atomic.AtomicLong;

public class SQLExeRecorder {
    private AtomicLong selectTime = new AtomicLong(0);
    private AtomicLong insertTime = new AtomicLong(0);
    private AtomicLong updateTime = new AtomicLong(0);

    private AtomicLong selectCount = new AtomicLong(0);
    private AtomicLong insertCount = new AtomicLong(0);
    private AtomicLong updateCount = new AtomicLong(0);

    private AtomicLong selectError = new AtomicLong(0);
    private AtomicLong insertError = new AtomicLong(0);
    private AtomicLong updateError = new AtomicLong(0);

    public void addSelectTime(long exeTime) {
        selectTime.addAndGet(exeTime);
        selectCount.addAndGet(1);
    }

    public void addSelectError() {
        selectError.addAndGet(1);
    }

    public void addInsertError() {
        insertError.addAndGet(1);
    }

    public void addUpdateError() {
        updateError.addAndGet(1);
    }

    public void addInsertTime(long exeTime) {
        insertTime.addAndGet(exeTime);
        insertCount.addAndGet(1);
    }

    public void addUpdateTime(long exeTime) {
        updateTime.addAndGet(exeTime);
        updateCount.addAndGet(1);
    }

    public void printSelectTimeStat() {
        try {
            System.out.println("selectTime:" + selectTime + ", selectCount:" + selectCount + ", selectError:" + selectError + ", avg:" + (selectTime.longValue() * 1.000 / selectCount.longValue()) + ", errorRate:" + (selectError.longValue() * 1.00000 / (selectCount.longValue() + selectError.longValue())));
        } catch (ArithmeticException e) {
            System.out.println("selectTime:" + selectTime + ", selectCount:" + selectCount);
        }
    }

    public void printInsertTimeStat() {
        try {
            System.out.println("insertTime:" + insertTime + ", insertCount:" + insertCount + ", insertError:" + insertError + ", avg:" + (insertTime.longValue() * 1.000 / insertCount.longValue()) + ", errorRate:" + (insertError.longValue() * 1.00000 / (insertError.longValue() + insertCount.longValue())));
        } catch (ArithmeticException e) {
            System.out.println("insertTime:" + insertTime + ", insertCount:" + insertCount);
        }
    }

    public void printUpdateTimeStat() {
        try {
            System.out.println("updateTime:" + updateTime + ", updateCount:" + updateCount + ", updateError:" + updateError + ", avg:" + (updateTime.longValue() * 1.000 / updateCount.longValue()) + ", errorRate:" + (updateError.longValue() * 1.00000 / (updateError.longValue() + updateCount.longValue())));
        } catch (ArithmeticException e) {
            System.out.println("updateTime:" + updateTime + ", updateCount:" + updateCount);
        }
    }


}