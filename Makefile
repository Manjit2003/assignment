build:
	@go build -o bin/api cmd/api/main.go 

swag:
	@swag init -g cmd/api/main.go

swag-tidy:
	@swag fmt -g cmd/api/main.go