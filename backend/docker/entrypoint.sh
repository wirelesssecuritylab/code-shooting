#!/bin/bash

mkdir -p /app/log

umask 127

exec /app/codeshooting >> /app/log/codeshooting.log 2>&1
