#!/bin/bash
port=9090
kill -9 $(lsof -i:$port |awk '{print $2}' | tail -n 2)
#./build/bin/bootnode -addr "$localip:$port" -nodekey=localBoot.key
./build/bin/bootnode -nodekey=localBoot.key
