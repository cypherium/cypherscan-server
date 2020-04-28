#!/bin/bash

localProjectPath="../cypherBFT/go-cypherium"
cp -rf $localProjectPath/crypto/bls/lib/mac/*     $localProjectPath/crypto/bls/lib/
cp -rf $localProjectPath/pow/cphash/randomX/lib/Darwin/*  $localProjectPath/pow/cphash/randomX/
go build -ldflags "-w" ./cmd/main/*
#cp -rf $localProjectPath/crypto/bls/lib/linux/*     $localProjectPath/crypto/bls/lib/
#cp -rf $localProjectPath/pow/cphash/randomX/lib/Linux/*   $localProjectPath/pow/cphash/randomX/

