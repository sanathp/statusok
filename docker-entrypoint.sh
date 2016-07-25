#!/bin/sh
set -e
sleep 5 # Workaround to wait untill InfluxDb will start
/go/bin/StatusOk --config /config/config.json
