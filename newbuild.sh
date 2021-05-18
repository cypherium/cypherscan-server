#!/bin/bash
rm -rf app
go build -o app cmd/main/*
./stop-server.sh
./start-server.sh
