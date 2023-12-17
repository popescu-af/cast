.PHONY: build-service
build-service: ## Build service binary
build-service:
	go build \
		-o ./build/cast-service \
		-tags osusergo,netgo \
		./app/service/...
		
.PHONY: build-ui
build-ui: ## Build service binary
build-ui:
	cd ui/cast_ui && flutter build web
	mkdir -p ui/dist && cp -r ui/cast_ui/build/web ui/dist/

.PHONY: up-service
up-service: ## Run only service
up-service: build-service
	docker compose up --build cast-service

.PHONY: down
down: ## Bring down docker containers
down:
	docker compose down

.PHONY: up
up: ## Bring up docker containers
up:
	docker compose up --build

.PHONY: run
run: ## Bring up docker containers in background
run:
	docker compose up --build -d
