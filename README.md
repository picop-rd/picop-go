# picop-go
PiCoP Protocol Library for Go

## Description
This repository provides Go libraries for communicating with PiCoP Protocol.

## Support
- `net`([`protocol/net`](./protocol/net`))
- `net/http`([`contrib/net/http/picophttp`](./contrib/net/http/picophttp))
- `github.com/go-sql-driver/mysql`([`contrib/github.com/go-sql-griver/mysql/picopmysql`](./contrib/github.com/go-sql-driver/mysql/picopmysql))
- `go.mongodb.org/mongo-driver/mongo`([`contrib/go.mongodb.org/mongo-driver/mongo/picopmongo`](./contrib/go.mongodb.org/mongo-driver/mongo/picopmongo))
  - Please use the forked version: [`github.com/picop-rd/mongo-go-driver`](https://github.com/picop-rd/mongo-go-driver)
- `github.com/bradfitz/gomemcache`([`contrib/github.com/bradfitz/gomemcache/picopgomemcache`](./contrib/github.com/bradfitz/gomemcache/picopgomemcache))
  - Please use the forked version: [`github.com/picop-rd/gomemcache`](https://github.com/picop-rd/gomemcache)
- [WIP] `google.golang.org/grpc`([`contrib/google.golang.org/grpc/picopgrpc`](./contrib/google.golang.org/grpc/picopgrpc))
  - Please use the forked version: [`github.com/picop-rd/grpc-go`](https://github.com/picop-rd/grpc-go)

You can refer to [`example`](./example) directory for usage.
