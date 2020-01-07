#!/usr/bin/env bash

npm i

node node_modules/@2dfire/server/bin/index build ENV=$1
