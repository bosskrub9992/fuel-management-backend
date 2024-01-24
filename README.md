### Getting Started (Manual approach)

1. start postgres using docker compose, run (recommend)
```sh
docker-compose up
```
or start postgres using only docker, run
```sh
docker run --rm -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e TZ=Asia/Bangkok --name postgres postgres:16-alpine
```

2. drop and create database, run (data will be lost!)
```sh
docker exec -i fuel-management-postgres psql -U postgres -c "drop database if exists fuel"
docker exec -i fuel-management-postgres psql -U postgres -c "create database fuel"
```

3. migrate database
```sh
go run ./cmd/up all
```

4. start service
```sh
go run main.go
```

### useful command

go environment for development
```sh
go env -w GONOPROXY="github.com/jinleejun-corp/*"
go env -w GONOSUMDB="github.com/jinleejun-corp/*"
go env -w GOPRIVATE="github.com/jinleejun-corp/*"
go env -w GOPROXY="https://proxy.golang.org,direct"
git config --global url."ssh://git@github.com".insteadOf "https://github.com"
```

to reset environment to default
```sh
go env -u GONOPROXY
go env -u GONOSUMDB
go env -u GOPRIVATE
go env -u GOPROXY
```

check vulnability
```sh
go install golang.org/x/vuln/cmd/govulncheck@latest
~/go/bin/govulncheck ./...
``` 