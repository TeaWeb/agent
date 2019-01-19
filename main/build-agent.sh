#!/usr/bin/env bash

# run service

export GOPATH=`pwd`/../../../../../

go build -o ${GOPATH}/src/main/bin/teaweb-agent ${GOPATH}/src/github.com/TeaWeb/agent/main/main-agent.go