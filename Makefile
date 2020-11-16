build-docker:
	docker build --tag sayitsocial:$(TAG)

run-docker:
	docker run --env-file .env sayitsocial:$(TAG)