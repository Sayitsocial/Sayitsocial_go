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
	docker build -t ghcr.io/sayitsocial/sayitsocial_go:development .

pull:
	docker pull ghcr.io/sayitsocial/sayitsocial_go:development && docker pull ghcr.io/sayitsocial/front-react:development

run:
	DEBUG=true go run main.go
