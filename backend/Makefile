COM := codeshooting

.PHONY: clean
clean:
	@rm -rf bin
	@echo Clean done!

.PHONY: build
build: BUILD_ENTRY = main.go
build: clean
	@rm -rf bin
	@echo Building...
	@mkdir -p bin/app
	@GOOS=linux GOARCH=$(ARCH) CGO_ENABLED=0 go build -o=bin/app/$(COM) -ldflags '-w -s' $(BUILD_ENTRY)
	@echo Build done!

.PHONY: docker-build
docker-build: build
	@cp docker/Dockerfile bin \
		&& cp docker/entrypoint.sh bin/app \
		&& cp -r conf bin/app \
		&& cd bin && chmod -R 755 app \
		&& DOCKER_BUILDKIT=1 docker build -t $(COM)-backend:$(VERSION) .
