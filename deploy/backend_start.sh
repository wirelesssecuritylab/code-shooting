#!/bin/bash

deploy_codeshooting() {
    docker run --name codeshooting-backend \
		-p 10000:2022 \
		-v /codeshooting/backend/data:/app/data \
		-v /codeshooting/backend/log:/app/log \
		-v /codeshooting/backend/conf:/app/conf \
		-v /etc/localtime:/etc/localtime:ro \
		-v /etc/timezone:/etc/timezone:ro \
		-d codeshooting-backend:latest
}

deploy_codeshooting
