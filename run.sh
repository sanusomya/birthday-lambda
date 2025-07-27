env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bootstrap main.go
zip bootstrap.zip bootstrap