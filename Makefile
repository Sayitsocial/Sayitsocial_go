include *.env
export

export TAG_GO=development
export TAG_REACT=development

uname_p := $(shell uname -p)
ifeq ($(uname_p),aarch64)
export TAG_PG=raspi
else
export TAG_PG=12-3.1-alpine
endif

up:
	docker-compose up -d

down:
	docker-compose down

build:
	docker-compose up --build

buildx-goapp:
	docker buildx build --push --platform linux/amd64,linux/arm64,linux/arm/v7 --tag sayitsocial/sayitsocial\:${TAG_GO} .

build-react:
	docker buildx build --push --platform linux/amd64,linux/arm64,linux/arm/v7 --tag sayitsocial/sayitsocial-react\:${TAG_REACT} ./web/v2/

push: buildx-goapp build-react

pull:
	docker pull sayitsocial/sayitsocial:development && docker pull sayitsocial/sayitsocial-react:development

run:
	DEBUG=true go run main.go
