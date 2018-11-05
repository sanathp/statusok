#!/bin/sh

cat /config.template | envsubst > /config.json
/statusok --config /config.json
