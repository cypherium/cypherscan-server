#!/usr/bin/env bash
systemctl stop server.service
go build ./cmd/*
systemctl start server.service
systemctl status server.service -l