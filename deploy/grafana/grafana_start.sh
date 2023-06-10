#!/bin/bash

deploy_grafana() {
	docker run --name codeshooting-grafana \
		-p 8899:3010 \
		-v /codeshooting/grafana/config/grafana.ini:/etc/grafana/grafana.ini \
		-v /codeshooting/grafana/data:/var/lib/grafana \
		-v /codeshooting/grafana/provisioning:/etc/grafana/provisioning \
		-v /etc/localtime:/etc/localtime:ro \
		-v /etc/timezone:/etc/timezone:ro \
		-idt grafana:8.4.5
}

deploy_grafana
