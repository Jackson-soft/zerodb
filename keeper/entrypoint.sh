#!/bin/sh

if [ -n "${ETC_ENV}" ]; then
    exec /opt/keeper/zeroKeeper -config /opt/keeper/etc/app_"${ETC_ENV}".yaml
else
    exec /opt/keeper/zeroKeeper -config /opt/keeper/etc/app_daily.yaml
fi
