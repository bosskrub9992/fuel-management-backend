### Getting Started (Manual approach)

1. start postgres, run
```sh
docker run --rm -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e TZ=Asia/Bangkok --name postgres postgres:16-alpine
```

2. drop and create database, run (data will be lost!)
```sh
docker exec -i postgres psql -U postgres -c "drop database if exists testdb" && \
docker exec -i postgres psql -U postgres -c "create database testdb"
```

3. start service
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