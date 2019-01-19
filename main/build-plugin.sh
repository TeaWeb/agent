#!/usr/bin/env bash

# run service

export GOPATH=`pwd`/../../../../../

go build -o ${GOPATH}/src/main/plugins/agent.tea ${GOPATH}/src/github.com/TeaWeb/agent/main/main-plugin.go