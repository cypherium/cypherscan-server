#!/usr/bin/env bash
cp -rf server.service  /etc/init.d/
chmod 700 /etc/init.d/server.service
rm -rf ./out*.log
systemctl daemon-reload
systemctl enable /etc/init.d/server.service
systemctl start server.service
systemctl status server.service -l
