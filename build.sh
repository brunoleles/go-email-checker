#!/bin/bash
cd `dirname $0`;

GOOS=windows GOARCH=amd64 go build -o ./bin/win-amd64/email-go-checker.exe

GOOS=linux go build -o ./bin/linux/email-go-checker