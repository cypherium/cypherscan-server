#!/bin/bash
killall -9 main
rm -rf ./out.log
localProjectPath="./vendor/github.com/cypherium/cypherBFT-P/go-cypherium"
cp -rf $localProjectPath/crypto/bls/lib/mac/*     $localProjectPath/crypto/bls/lib/
cp -rf $localProjectPath/pow/cphash/randomX/lib/Darwin/*  $localProjectPath/pow/cphash/randomX/
go build -o scansv ./cmd/*
cp -rf $localProjectPath/crypto/bls/lib/linux/*     $localProjectPath/crypto/bls/lib/
cp -rf $localProjectPath/pow/cphash/randomX/lib/Linux/*   $localProjectPath/pow/cphash/randomX/
./scansv >>./out.log 2>&1
