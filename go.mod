module github.com/go-trellis/trellis

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0

go 1.13

require (
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-trellis/common v1.7.0
	github.com/go-trellis/config v1.4.1
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/mailru/easyjson v0.7.1
)
