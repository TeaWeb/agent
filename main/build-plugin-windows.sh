#!/usr/bin/env bash

# run service

export GOPATH=`pwd`/../../../../../
export GOOS=windows

go build -o ${GOPATH}/src/main/plugins/agent.tea.exe ${GOPATH}/src/github.com/TeaWeb/agent/main/main-plugin.go