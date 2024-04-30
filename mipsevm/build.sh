#!/usr/bin/env bash
set -e

cd ../minigeth
export GOOS=linux
export GOARCH=mips64
export GOMIPS=softfloat
go build -o ../mipsevm/minigeth

cd ../mipsevm
file minigeth

