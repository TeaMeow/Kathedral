rm ./bin/main
go build -o ./bin/main
cd ./bin
./main --token INSERT_THE_TOKEN_HERE

# For Windows
# GOOS=windows GOARCH=386 go build -o kathedral.exe

# For Linux
# GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o kathedral.linux