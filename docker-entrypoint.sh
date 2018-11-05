#!/bin/sh

cat /config.template | envsubst > /config.json
cat /config.json
/statusok --config /config.json
