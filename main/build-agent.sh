#!/usr/bin/env bash

# run service

export GOPATH=`pwd`/../../../../../
#export GOOS=linux
#export GOARCH=amd64

go build -o ${GOPATH}/src/main/bin/teaweb-agent ${GOPATH}/src/github.com/TeaWeb/agent/main/main-agent.go