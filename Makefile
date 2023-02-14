.PHONY: docker-up
docker-up:
	docker-compose -f docker-compose.yaml up --build

.PHONY: docker-down
docker-down: ## Stop docker containers and clear artefacts.
	docker-compose -f docker-compose.yaml down
	docker system prune 

.PHONY: bundle
bundle: ## bundles the submission for... submission
	git bundle create guestlist.bundle --all

run-api:
	@go run cmd/app/main.go
.PHONY: run-api