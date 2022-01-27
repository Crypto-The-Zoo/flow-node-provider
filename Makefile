.PHONY: clean test security build run

APP_NAME = apiserver
BUILD_DIR = $(PWD)/build
MIGRATIONS_FOLDER = $(PWD)/platform/migrations
DATABASE_URL = postgres://postgres:password@localhost/inception?sslmode=disable
STAGING_DATABASE_PROXY_URL = postgres://postgres:ND88Kdssu9dDE0lo@localhost:5544/inception-staging?sslmode=disable
PRODUCTION_DATABASE_PROXY_URL = postgres://postgres:ND88Kdssu9dDE0lo@localhost:5544/inception-production?sslmode=disable

clean:
	rm -rf ./build

security:
	gosec -quiet ./...

test: security
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

build: clean test
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

run: swag build
	$(BUILD_DIR)/$(APP_NAME)

migrate.up:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" up

migrate.down:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" down

migrate.force:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" force $(version)

docker.run: docker.network docker.postgres

docker.network:
	docker network inspect dev-network >/dev/null 2>&1 || \
	docker network create -d bridge dev-network

docker.fiber.build:
	docker build -t fiber .

docker.fiber: docker.fiber.build
	docker run --rm -d \
		--name dev-fiber \
		--network dev-network \
		-p 5000:5000 \
		fiber

docker.postgres:
	docker run --rm -d \
		--name dev-postgres \
		--network dev-network \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=password \
		-e POSTGRES_DB=postgres \
		-v ${HOME}/dev-postgres/data/:/var/lib/postgresql/data \
		-p 5432:5432 \
		postgres

docker.stop: docker.stop.fiber docker.stop.postgres

docker.stop.fiber:
	docker stop dev-fiber

docker.stop.postgres:
	docker stop dev-postgres

swag:
	swag init

deploy.staging:
	gcloud beta run deploy api-staging --set-env-vars='ENV=staging' --port=80 --add-cloudsql-instances=crypto-the-zoo-staging:us-west1:inception-db --project crypto-the-zoo-staging

db.socket:
	./bin/cloud_sql_proxy -instances=crypto-the-zoo-staging:us-west1:inception-db=tcp:5544

migrate.staging.up:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(STAGING_DATABASE_PROXY_URL)" up
	
migrate.staging.down:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(STAGING_DATABASE_PROXY_URL)" down

migrate.production.up:	
	migrate -path $(MIGRATIONS_FOLDER) -database "$(PRODUCTION_DATABASE_PROXY_URL)" up
	
migrate.production.down:
	migrate -path $(MIGRATIONS_FOLDER) -database "$(PRODUCTION_DATABASE_PROXY_URL)" down

deploy.production:
	gcloud beta run deploy api-production --set-env-vars='ENV=prod' --port=80 --add-cloudsql-instances=crypto-the-zoo-staging:us-west1:inception-db --project crypto-the-zoo-staging