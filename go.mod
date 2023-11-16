module github.com/picop-rd/picop-go

go 1.18

require go.opentelemetry.io/otel v1.11.2

require (
	github.com/bradfitz/gomemcache v0.0.0-20230905024940-24af94b03874
	github.com/go-sql-driver/mysql v1.7.0
	github.com/google/go-cmp v0.6.0
	go.mongodb.org/mongo-driver v1.12.1
	golang.org/x/sync v0.5.0
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231106174013-bbf56f31fb17 // indirect
)

require (
	github.com/golang/snappy v0.0.1 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.opentelemetry.io/otel/trace v1.11.2 // indirect
	golang.org/x/crypto v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/grpc v1.59.0
)

replace (
	github.com/bradfitz/gomemcache => github.com/picop-rd/gomemcache v1.0.0-picop
	go.mongodb.org/mongo-driver => github.com/picop-rd/mongo-go-driver v1.12.1-picop
	google.golang.org/grpc => github.com/picop-rd/grpc-go v1.0.1-picop
)
