module Worker

go 1.18

require google.golang.org/protobuf v1.26.0

require (
	Lib/db v0.0.0-00010101000000-000000000000 // indirect
	Lib/models v0.0.0-00010101000000-000000000000 // indirect
	Lib/rediswrapper v0.0.0-00010101000000-000000000000 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.mongodb.org/mongo-driver v1.11.1 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/text v0.3.7 // indirect
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fabrizio2210/photobook v0.0.0-00010101000000-000000000000
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
)

replace (
	Lib/db => ../lib/db
	Lib/models => ../lib/models
	Lib/rediswrapper => ../lib/rediswrapper
	github.com/fabrizio2210/photobook => ../lib/github.com/fabrizio2210/photobook
)
