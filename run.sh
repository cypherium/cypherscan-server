#!/bin/bash
killall -9 serverbackend
ps -ef|grep 8000|grep -v grep|cut -c 9-15|xargs kill -9
rm -rf ./out.log
localProjectPath="./vendor/github.com/cypherium/cypherBFT"
cp -rf $localProjectPath/crypto/bls/lib/darwin/*     $localProjectPath/crypto/bls/lib/
go build -o serverbackend ./cmd/*
cp -rf $localProjectPath/crypto/bls/lib/linux/*     $localProjectPath/crypto/bls/lib/
nohup ./serverbackend >>./out.log 2>&1 &
