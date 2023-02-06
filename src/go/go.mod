module fabrizio2210/Worker

go 1.18

require google.golang.org/protobuf v1.26.0

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
        github.com/fabrizio2210/photobook v0.0.0-00010101000000-000000000000
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
)

replace (
  github.com/fabrizio2210/photobook => ./github.com/fabrizio2210/photobook
)
