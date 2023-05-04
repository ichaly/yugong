CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main/dist/app main/main.go

docker build -t yugong:latest .

rm -rf main/dist