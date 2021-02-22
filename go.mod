module github.com/iTrellis/trellis

go 1.13

replace (
	go.etcd.io/etcd/api/v3 v3.5.0-pre => go.etcd.io/etcd/api/v3 v3.0.0-20210107172604-c632042bb96c
	go.etcd.io/etcd/pkg/v3 v3.5.0-pre => go.etcd.io/etcd/pkg/v3 v3.0.0-20210107172604-c632042bb96c
)

require (
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/pprof v1.3.0
	github.com/gin-gonic/gin v1.6.3
	github.com/golang/protobuf v1.4.3
	github.com/golang/snappy v0.0.2 // indirect
	github.com/google/uuid v1.2.0
	github.com/iTrellis/common v0.21.1
	github.com/iTrellis/config v0.21.1
	github.com/iTrellis/node v0.21.1
	github.com/iTrellis/xorm_ext v0.21.2
	github.com/mitchellh/hashstructure/v2 v2.0.1
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c
	github.com/urfave/cli/v2 v2.3.0
	go.etcd.io/etcd/api/v3 v3.5.0-pre
	go.etcd.io/etcd/client/v3 v3.0.0-20210201223203-e897daaebc2f
	go.uber.org/zap v1.16.0
	google.golang.org/grpc v1.29.1
	xorm.io/xorm v1.0.7
)
