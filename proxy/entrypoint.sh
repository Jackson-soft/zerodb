#!/bin/sh

if [ -n "${ETC_ENV}" ]; then
    exec /opt/proxy/zeroProxy -config /opt/proxy/etc/app_"${ETC_ENV}".yaml
else
    exec /opt/proxy/zeroProxy -config /opt/proxy/etc/app_daily.yaml
fi
