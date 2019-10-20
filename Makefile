test:
	go test ./compact

generate:
	go run cmd/generateforms/main.go
	gofmt -w compact/*
