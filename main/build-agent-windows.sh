#!/usr/bin/env bash

# run service

export GOPATH=`pwd`/../../../../../
export GOOS=windows

go build -o ${GOPATH}/src/main/bin/teaweb-agent.exe ${GOPATH}/src/github.com/TeaWeb/agent/main/main-agent.go