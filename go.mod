module github.com/go-trellis/trellis

// replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0

go 1.13

require (
	github.com/coreos/etcd v3.3.19+incompatible // indirect
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-resty/resty/v2 v2.3.0
	github.com/go-trellis/cache v1.1.1
	github.com/go-trellis/common v1.8.2
	github.com/go-trellis/config v1.4.6
	github.com/go-trellis/node v1.2.2
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/spf13/cobra v1.0.0
	go.etcd.io/etcd v3.3.19+incompatible
	go.uber.org/zap v1.15.0 // indirect
	google.golang.org/grpc v1.26.0
	sigs.k8s.io/yaml v1.2.0 // indirect
)
