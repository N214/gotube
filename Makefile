# Get the Git commit hash
GIT_COMMIT_HASH := $(shell git rev-parse HEAD)

# Define the image name
IMAGE_NAME := northamerica-northeast1-docker.pkg.dev/$(PROJECT)/cloud-run-source-deploy/gotube:$(GIT_COMMIT_HASH)
PORT := 8080
SERVICE_NAME := yt-notification-bot

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: check-project

check-project:
	@if [ -z "$(PROJECT)" ]; then \
		echo "ERROR: 'PROJECT' environment variable is not set."; \
		exit 1; \
	fi

build: check-project
	@docker build -t $(IMAGE_NAME) .

push: build
	@docker push $(IMAGE_NAME)
	@echo "Container pushed to $(IMAGE_NAME) in Artifact Registry."

deploy: push
	@gcloud run deploy $(SERVICE_NAME) \
		--image=$(IMAGE_NAME) \
		--execution-environment=gen2 \
		--platform=managed \
		--region=northamerica-northeast1 \
		--port=$(PORT) \
		--allow-unauthenticated
        # Add other Cloud Run deployment options here
	@echo "Container deployed to Cloud Run."