BINARY = ctfzone

pkgs = $(shell go list ./...)

build:
	@echo ">> building binaries"
	@go build -ldflags "-s -w" -o ${BINARY}

clean:
	@echo ">> cleaning binaries"
	@if [ -f ${BINARY} ]; then rm ${BINARY}; fi

format:
	@echo ">> formatting code"
	@go fmt ${pkgs}

test:
	@echo ">> running all tests"
	@go test -p 1 ${pkgs}

tools:
	@echo ">> install nesessary tools"
	@go get -u github.com/jteeuwen/go-bindata/...
	@go get -u github.com/vektra/mockery/...

pack-schemas:
	@echo ">> writing json-schemas"
	@go-bindata -pkg schemas \
		-ignore ".*\.go" \
		-prefix "controllers/schemas" \
		-o controllers/schemas/bindata.go \
		controllers/schemas/...
	@go fmt ./controllers/schemas

pack-migrations:
	@echo ">> writing migrations"
	@go-bindata -pkg migrations \
		-ignore ".*\.go" \
		-prefix "models/migrations" \
		-o models/migrations/bindata.go \
		models/migrations/...
	@go fmt ./models/migrations

mailer-mock:
	@echo ">> generating mailer mock"
	@mockery -dir modules/mailer \
		-output modules/mailer/mock \
		-outpkg mailer_mock \
		-name Sender
	@go fmt ./modules/mailer/mock

db-clean:
	@echo ">> cleaning database"
	@psql ${CTF_DB_DSN} -c "DROP OWNED BY ctfzone"

migrations-up:
	@echo ">> migrations up"
	@migrate -database ${CTF_DB_DSN} -path ./models/migrations up

migrations-down:
	@echo ">> migrations down"
	@migrate -database ${CTF_DB_DSN} -path ./models/migrations down
