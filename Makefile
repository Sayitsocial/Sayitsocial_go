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

pull:
	docker pull sayitsocial/sayitsocial:development && docker pull sayitsocial/sayitsocial-react:development

run:
	DEBUG=true go run main.go
