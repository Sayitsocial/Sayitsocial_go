include *.env
export

export TAG_GO=development
export TAG_REACT=development

up:
	docker-compose up -d

down:
	docker-compose down

build:
	docker-compose up --build

buildx-goapp:
	docker buildx build --push --platform linux/amd64,linux/arm64,linux/arm/v7 --tag ${REGISTRY_PATH_GO}/sayitsocial\:${TAG_GO} .

build-react:
	docker buildx build --push --platform linux/amd64,linux/arm64,linux/arm/v7 --tag ${REGISTRY_PATH_REACT}/sayitsocial-react\:${TAG_REACT} ./web/v2/
