include .env
ifneq ("$(wildcard .env.local)","")
	include .env.local
endif
export

CMD=go run ./cmd/test-task-manager/main.go -loglevel=debug

run:
	${CMD}

migrate:
	${CMD} -migrate

migrate-down:
	${CMD} -migrate down

test:
	go test -v -coverprofile=coverage.out `go list ./internal/... ./pkg/... | grep -Ev "/app|/http/server"` && go tool cover -html=coverage.out
