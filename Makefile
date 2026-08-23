APP = sam2ooh2

run:
	@go run cmd/sam2ooh2/sam2ooh2.go

package:
	@go mod tidy

migrate: 
	@go run cmd/migrate/migrate.go