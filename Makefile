include .env
export

export UID := $(shell id -u)
export GID := $(shell id -g)
export PROJECT_ROOT=${shell pwd}

env-up:
	@mkdir -p out/pgdata
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres

env-cleanup:
	@read -p "Очистить все volume файлы окружения? Опасность утери данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todoapp-postgres port-forwarder && \
		sudo rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose stop port-forwarder

migrate-create:
	@if [ -z "${seq}" ]; then \
		echo "Отсутствует необходимый параметр seq. Пример make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		--user $(shell id -u):$(shell id -g) \
		create \
		-ext sql \
		-dir /migrations \
		-seq "${seq}"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "${action}" ]; then \
		echo "Отсутствует необходимый параметр action. Пример make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"${action}"

migrate-force:
	@if [ -z "${version}" ]; then \
		echo "Отсутствует необходимый параметр version. Пример make migrate-force version=1"; \
		exit 1; \
	fi; \
	docker compose run --rm \
		todoapp-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		force ${version}

fix-permissions:
	@sudo chown -R $(shell id -u):$(shell id -g) out/ migrations/

logs-cleanup:
	@read -p "Очистить все лог файлы окружения? Опасность утери лог данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		sudo rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "Файлы логов очищены"; \
	else \
		echo "Очистка логов отменена"; \
	fi

todoapp-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/todoapp/main.go

todoapp-deploy:
	docker compose up -d --build todoapp

todoapp-undeploy:
	@docker compose down todoapp

swagger-gen:
	@docker compose run --rm swagger \
		init \
		-g cmd/todoapp/main.go \
		-o docs \
		--parseInternal \
		--parseDependency

ps:
	@docker compose ps