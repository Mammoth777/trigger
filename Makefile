linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./release/trigger-server-linux main.go
mac:
	go build -o ./release/trigger-server main.go