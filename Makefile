prepare:
	@echo "Preparing..."
	@cd libs && go mod download

dev.test:
	@echo "Run tests..."
	@cd libs && go test -v ./...

dev.infra.up:
	@echo "Starting development infrastructure..."
	@docker-compose -f docker-compose.infra.yml up -d

dev.migrate.up:
	@echo "Starting development migration..."
	@docker-compose -f docker-compose.migrate.yml up

dev.up.watch:
	@echo "Starting development services..."
	@docker-compose -f docker-compose.yml up

dev.up:
	@echo "Starting development services..."
	@docker-compose -f docker-compose.yml up -d

dev.down:
	@echo "Starting development services..."
	@docker-compose -f docker-compose.yml down

dev.down.all:
	@echo "Starting development services..."
	@docker-compose -f docker-compose.yml down
	@docker-compose -f docker-compose.infra.yml down
	@docker-compose -f docker-compose.migrate.yml down
