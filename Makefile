APP = sam2ooh2

init:
	@go mod tidy
	@mockery

run:
	@go run cmd/sam2ooh2/sam2ooh2.go

package:
	@go mod tidy

migrate: 
	@go run cmd/migrate/migrate.go