package etcd

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	bsf "github.com/iTrellis/common/encryption/binary-formats"
	"github.com/iTrellis/common/errors"
	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/registry"

	"github.com/google/uuid"
	"github.com/mitchellh/hashstructure/v2"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type etcdRegistry struct {
	id      string
	options registry.Options

	// logger logger.Logger

	sync.RWMutex

	services map[string]register
	leases   map[string]leases

	workers map[string]workers

	client *clientv3.Client
}

type register map[string]uint64
type leases map[string]clientv3.LeaseID
type workers map[string]*worker

// NewRegistry new etcd registry
func NewRegistry(opts ...registry.Option) (registry.Registry, error) {

	p := &etcdRegistry{
		id: uuid.New().String(),

		// map[registryFullPath]map[node.ID]id
		services: make(map[string]register),
		// map[registryFullPath]map[node.ID]leaseID
		leases: make(map[string]leases),
		// map[registryFullPath]map[node.ID]worker
		workers: make(map[string]workers),
	}

	configure(p, opts...)

	return p, nil
}

// configure will setup the registry with new options
func configure(e *etcdRegistry, opts ...registry.Option) error {
	for _, o := range opts {
		o(&e.options)
	}

	// setup the client
	cli, err := newClient(e)
	if err != nil {
		return err
	}

	if e.client != nil {
		e.client.Close()
	}

	// setup new client
	e.client = cli

	return nil
}

func newClient(e *etcdRegistry, opts ...registry.Option) (*clientv3.Client, error) {
	for _, o := range opts {
		o(&e.options)
	}

	config := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	var cAddrs []string
	for _, address := range e.options.Addrs {
		if len(address) == 0 {
			continue
		}
		addr, port, err := net.SplitHostPort(address)
		if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
			port = "2379"
			addr = address
			cAddrs = append(cAddrs, net.JoinHostPort(addr, port))
		} else if err == nil {
			cAddrs = append(cAddrs, net.JoinHostPort(addr, port))
		}
	}

	// if we got addrs then we'll update
	if len(cAddrs) > 0 {
		config.Endpoints = cAddrs
	}

	if e.options.Timeout != 0 {
		config.DialTimeout = e.options.Timeout
	}

	if e.options.Secure || e.options.TLSConfig != nil {
		tlsConfig := e.options.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		config.TLS = tlsConfig

		for i, ep := range config.Endpoints {
			if !strings.HasPrefix(ep, "https://") {
				config.Endpoints[i] = "https://" + ep
			}
		}
	}

	if e.options.Context != nil {
		u, ok := e.options.Context.Value(authKey{}).(*authCreds)
		if ok {
			config.Username = u.Username
			config.Password = u.Password
		}
		cfg, ok := e.options.Context.Value(logConfigKey{}).(*zap.Config)
		if ok && cfg != nil {
			config.LogConfig = cfg
		}
	}

	// e.logger.Info("etcd config", config)

	return clientv3.New(config)
}

func (p *etcdRegistry) ID() string {
	return p.id
}

func (p *etcdRegistry) Init(opts ...registry.Option) error {
	return configure(p, opts...)
}

func (p *etcdRegistry) Options() registry.Options {
	return p.options
}

func (p *etcdRegistry) Register(s *registry.Service, opts ...registry.RegisterOption) error {
	if s.Node == nil {
		return errors.New("Require node's infor")
	}

	var options registry.RegisterOptions
	for _, o := range opts {
		o(&options)
	}

	if options.Heartbeat == 0 {
		options.Heartbeat = 10 * time.Second
	}

	// p.logger.Debug("etcd_registry", "register_service", s.Service, "options", options)

	p.RLock()
	wors, ok := p.workers[s.Service.FullRegistry()]
	p.RUnlock()

	if !ok {
		wors = make(workers)
	}

	nodePath := s.Service.FullRegistry(s.Node.Value)
	for _, w := range wors {
		if nodePath == w.fullpath || w.service.Node.ID == s.Node.ID {
			// no need register again
			return nil
		}
	}

	wer := &worker{
		service:    s,
		ticker:     time.NewTicker(options.Heartbeat),
		fullpath:   nodePath,
		stopSignal: make(chan bool),
		options:    options,
	}

	go func(wr *worker) {
		var count uint32
		for {
			if err := p.registerNode(s, wr); err != nil {
				if wr.options.RetryTimes == 0 {
					continue
				}
				if wr.options.RetryTimes <= count {
					panic(fmt.Errorf("%s regist into etcd failed times above: %d, %v", wr.fullpath, count, err))
				}
				count++
				continue
			}

			count = 0
			select {
			case <-wr.stopSignal:
				p.Deregister(s)
				return
			case <-wr.ticker.C:
				// nothing to do
			}
		}
	}(wer)

	p.Lock()
	wors[s.Node.ID] = wer
	p.workers[s.Service.FullRegistry()] = wors
	p.Unlock()

	return nil
}

func (p *etcdRegistry) registerNode(s *registry.Service, wr *worker) error {
	if s == nil || s.Node == nil || s.Node.ID == "" || s.Node.Value == "" {
		return errors.New("node should not be nil")
	}
	regFullpath := s.FullRegistry()

	p.Lock()
	// ensure the leases and registers are setup for this domain
	if _, ok := p.leases[regFullpath]; !ok {
		p.leases[regFullpath] = make(leases)
	}
	if _, ok := p.services[regFullpath]; !ok {
		p.services[regFullpath] = make(register)
	}

	leaseID, ok := p.leases[regFullpath][s.Value]
	p.Unlock()

	// p.logger.Debug("register service node", ok, *s)

	fullregistryPath := s.FullRegistry(s.Node.ID)
	if !ok {
		// minimum lease TTL is ttl-second
		ctx, cancel := context.WithTimeout(context.Background(), p.options.Timeout)
		defer cancel()
		resp, err := p.client.Get(ctx, fullregistryPath, clientv3.WithSerializable())
		if err != nil {
			return err
		}
		for _, kv := range resp.Kvs {
			if kv.Lease <= 0 {
				continue
			}
			leaseID = clientv3.LeaseID(kv.Lease)

			// decode the existing node
			srv := decode(kv.Value)
			if srv == nil || srv.Node == nil {
				continue
			}

			h, err := hashstructure.Hash(srv, hashstructure.FormatV2, nil)
			if err != nil {
				return err
			}

			// save the info
			p.Lock()
			p.leases[regFullpath][s.Node.ID] = leaseID
			p.services[regFullpath][s.Node.ID] = h
			p.Unlock()

			break
		}
	}

	var leaseNotFound bool
	if leaseID > 0 {
		if _, err := p.client.KeepAliveOnce(context.TODO(), leaseID); err != nil {
			if err != rpctypes.ErrLeaseNotFound {
				return err
			}
		}

		leaseNotFound = true
	}

	// create hash of service; uint64
	h, err := hashstructure.Hash(s.Node, hashstructure.FormatV2, nil)
	if err != nil {
		return err
	}

	// get existing hash for the service node
	p.RLock()
	v, ok := p.services[regFullpath][s.Node.ID]
	p.RUnlock()

	if ok && v == h && !leaseNotFound {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.options.Timeout)
	defer cancel()

	var lgr *clientv3.LeaseGrantResponse
	if wr.options.TTL.Seconds() > 0 {
		// get a lease used to expire keys since we have a ttl
		lgr, err = p.client.Grant(ctx, int64(wr.options.TTL.Seconds()))
		if err != nil {
			return err
		}
	}
	var putOpts []clientv3.OpOption
	if lgr != nil {
		putOpts = append(putOpts, clientv3.WithLease(lgr.ID))
	}

	if _, err = p.client.Put(ctx, fullregistryPath, encode(s), putOpts...); err != nil {
		return err
	}

	p.Lock()
	// save our hash of the service
	p.services[regFullpath][s.Node.ID] = h
	// save our leaseID of the service
	if lgr != nil {
		p.leases[regFullpath][s.Node.ID] = lgr.ID
	}
	p.Unlock()

	return nil
}

func (p *etcdRegistry) Deregister(s *registry.Service, opts ...registry.DeregisterOption) error {
	if s == nil || s.Node == nil || s.Node.ID == "" {
		return errors.New("node should not be nil")
	}

	regServicePath := s.FullRegistry()

	p.Lock()
	// delete our hash of the service
	nodes, ok := p.services[regServicePath]
	if ok {
		delete(nodes, s.Node.ID)
		p.services[regServicePath] = nodes
	}

	// delete our lease of the service
	leases, ok := p.leases[regServicePath]
	if ok {
		delete(leases, s.Node.ID)
		p.leases[regServicePath] = leases
	}

	workers, ok := p.workers[regServicePath]
	if ok {
		delete(workers, s.Node.ID)
	}

	p.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), p.options.Timeout)
	defer cancel()
	if _, err := p.client.Delete(ctx, s.FullRegistry(s.Node.ID)); err != nil {
		return err
	}

	return nil
}

func (p *etcdRegistry) String() string {
	return service.RegisterType_name[int32(service.RegisterType_etcd)]
}

func (p *etcdRegistry) Stop() error {
	p.Lock()
	p.Unlock()
	for _, workers := range p.workers {
		for _, w := range workers {
			w.stopSignal <- true

			ctx, cancel := context.WithTimeout(context.Background(), p.options.Timeout)
			defer cancel()
			if _, err := p.client.Delete(ctx, w.service.FullRegistry(w.service.Node.ID)); err != nil {
				return err
			}
		}
	}

	if p.client != nil {
		p.client.Close()
	}

	return nil
}

func (p *etcdRegistry) Watch(opts ...registry.WatchOption) (registry.Watcher, error) {
	cli, err := newClient(p)
	if err != nil {
		return nil, err
	}
	return newEtcdWatcher(cli, p.id, p.options.Timeout, opts...)
}

func encode(nn *registry.Service) string {
	bs, _ := json.Marshal(nn)
	return bsf.Encode(bsf.EncodeStd, bs)
}

func decode(bs []byte) *registry.Service {
	dst, err := bsf.Decode(bsf.EncodeStd, bs)
	if err != nil {
		return nil
	}

	var s *registry.Service
	json.Unmarshal(dst, &s)
	return s
}
