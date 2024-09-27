build:
	@go build -o bin/personal-backend

run: build
	@./bin/personal-backend
test:
	@go test -b ./..
