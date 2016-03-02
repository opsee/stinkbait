PROJECT=stinkbait

all: fmt build

build:
	govendor build cmd/$(PROJECT)

fmt:
	@gofmt -w ./

deps:
	docker-compose up -d

docker: deps fmt dbuild

run: docker drun

dbuild:
	docker run \
		--link $(PROJECT)_memcached:memcached \
		--env-file ./$(APPENV) \
		-e "TARGETS=linux/amd64" \
		-e GODEBUG=netdns=cgo \
		-e PROJECT=github.com/opsee/$(PROJECT) \
		-v `pwd`:/gopath/src/github.com/opsee/$(PROJECT) quay.io/opsee/build-go:16 \
		&& docker build -t quay.io/opsee/$(PROJECT) .

drun:
	docker run \
		--link $(PROJECT)_memcached:memcached \
		--env-file ./$(APPENV) \
		-e GODEBUG=netdns=cgo \
		-e AWS_DEFAULT_REGION \
		-e AWS_ACCESS_KEY_ID \
		-e AWS_SECRET_ACCESS_KEY \
		-p 9100:9100 \
		--rm \
		quay.io/opsee/$(PROJECT):latest

.PHONY: docker dbuild drun run migrate clean all
