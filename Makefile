test:
	go test .

generate:
	go run cmd/generateforms/main.go
	gofmt -w .
