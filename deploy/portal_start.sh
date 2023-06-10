#!/bin/bash

deploy_codeshootingportal() {
	docker run --name codeshooting-portal \
		-p 80:80 \
		-p 443:443 \
		-v /codeshooting/portal/nginx/conf.d/default.conf:/etc/nginx/conf.d/default.conf \
		-v /etc/localtime:/etc/localtime:ro \
		-v /etc/timezone:/etc/timezone:ro \
		-d codeshooting-portal:latest
}

deploy_codeshootingportal
