# Makefile for migrations and dev deps
# expects a .env file with DB_ADDR (or export DB_ADDR in your environment)
-include .env
export DB_ADDR

MIGRATION_PATH := ./migrate/migrations
MIGRATE := migrate

.PHONY: dep-up dep-down create-migration migrate-up migrate-down \
        migrate-force migrate-version seed

dep-up:
	@docker compose up -d

dep-down:
	@docker compose down

# usage:
#   make create-migration name_here
#   e.g. make create-migration add_followers_table
create-migration:
	@$(MIGRATE) create -seq -ext sql -dir $(MIGRATION_PATH) $(filter-out $@,$(MAKECMDGOALS))

# run all up migrations
migrate-up:
	@$(MIGRATE) -path=$(MIGRATION_PATH) -database="$(DB_ADDR)" up

# run down. Accepts either:
#   make migrate-down 9        # passing a bare goal/arg like you were using
#   make migrate-down VERSION=9
#   make migrate-down          # runs one "down" step (depends on migrate binary)
migrate-down:
	# prefer explicit VERSION if provided, otherwise forward any extra goals
	@if [ -n "$(VERSION)" ]; then \
	  $(MIGRATE) -path=$(MIGRATION_PATH) -database="$(DB_ADDR)" down $(VERSION); \
	else \
	  $(MIGRATE) -path=$(MIGRATION_PATH) -database="$(DB_ADDR)" down $(filter-out $@,$(MAKECMDGOALS)); \
	fi

# force the recorded schema version (clears dirty flag) -- use with extreme caution
# usage: make migrate-force VERSION=9
migrate-force:
	@if [ -z "$(VERSION)" ]; then \
	  echo "Usage: make migrate-force VERSION=<version>"; exit 1; \
	fi
	@$(MIGRATE) -path=$(MIGRATION_PATH) -database="$(DB_ADDR)" force $(VERSION)

# query current recorded version
migrate-version:
	@$(MIGRATE) -path=$(MIGRATION_PATH) -database="$(DB_ADDR)" version

seed:
	@go run migrate/seed/main.go
