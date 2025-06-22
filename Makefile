#
# INTERNAL VARIABLES
#

#
# TARGETS
#
create-migration:
	@goose -dir=migrations/ create $(name)

run:
	@echo "[run] Running service in debug-hot-reload mode..."
	@export $$(cat dev.env) && nodemon --exec go run cmd/server/main.go --signal SIGTERM

run-simple:
	@echo "[run-simple] Running service..."
	@export $$(cat dev.env) && go run cmd/server/main.go

migrate:
	@echo "[migrate] Running database migrations..."
	@export $$(cat dev.env) && goose -dir=migrations/ -driver=sqlite3 -dsn=./data/admoai.db up

clean:
	@echo "[clean] Cleaning up database files..."
	@rm -rf ./data

setup:
	@echo "[setup] Setting up the project..."
	@go mod tidy
	@mkdir -p ./data
	@make migrate



