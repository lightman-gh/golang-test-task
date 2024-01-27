# Run compose
```
docker-compose up -d
```

# Run task 1
```
go build -o cmd_api ./cmd/cmd_api/main.go
./cmd_api
```

# Run task 2
```
go build -o cmd_processor ./cmd/cmd_processor/main.go
./cmd_processor
```

# Run task 3
```
go build -o cmd_reporting ./cmd/cmd_reporting/main.go
./cmd_reporting
```