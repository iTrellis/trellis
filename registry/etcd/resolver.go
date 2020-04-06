package etcd

import (
	"context"
	"time"

	"github.com/go-trellis/trellis/registry"

	"github.com/go-trellis/etcdnaming"
)

// NewResolver 生成对象
func NewResolver(opts ...registry.ResolverOptionFunc) {
	opt := &registry.ResolverOptions{}
	for _, o := range opts {
		o(opt)
	}
	if opt.Timeout == 0 {
		opt.Timeout = 8 * time.Second
	}
	etcdnaming.NewBuilder(etcdnaming.BuilderOptions{

		Server:     opt.GetName(),
		Version:    opt.GetVersion(),
		Endpoint:   opt.Target,
		LooperTime: opt.Timeout,
	})
}

func DialContext(ctx context.Context) {

}
