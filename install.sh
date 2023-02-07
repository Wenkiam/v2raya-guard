#!/bin/bash
go build
WORK_DIR=/usr/local/v2raya
if [ ! -d "$WORK_DIR" ];then
  mkdir -p /usr/local/v2raya
fi
mv v2raya-guard $WORK_DIR/
mv config.json $WORK_DIR/
mv v2raya-guard.service /usr/lib/systemd/system/

systemctl enable v2raya-guard
systemctl start v2raya-guard