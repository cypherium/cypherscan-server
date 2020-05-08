#!/usr/bin/env bash
cp server.service  /etc/init.d/
chmod 700 /etc/init.d/server.service

systemctl daemon-reload
systemctl enable /etc/init.d/server.service
systemctl start server.service
systemctl status server.service -l