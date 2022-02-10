#!/bin/bash
killall -9 server
rm -rf ./out.log
localProjectPath="./vendor/github.com/cypherium/cypherBFT"
cp -rf $localProjectPath/crypto/bls/lib/mac/*     $localProjectPath/crypto/bls/lib/
cp -rf $localProjectPath/pow/cphash/randomX/lib/Darwin/*  $localProjectPath/pow/cphash/randomX/
go build -o server ./cmd/*
cp -rf $localProjectPath/crypto/bls/lib/linux/*     $localProjectPath/crypto/bls/lib/
cp -rf $localProjectPath/pow/cphash/randomX/lib/Linux/*   $localProjectPath/pow/cphash/randomX/
nohup ./server >>./out.log 2>&1 &
