package routes

import (
	"context"
	"fmt"

	"github.com/iTrellis/node"
	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/client/grpc"
	"github.com/iTrellis/trellis/service/component"
	"github.com/iTrellis/trellis/service/message"
	"github.com/iTrellis/trellis/service/router"
)

// NewManager routes manager
func NewManager(opts ...Option) Manager {

	r := &manager{}

	r.Init(opts...)

	return r
}

type manager struct {
	// Router      router.Router
	// CompManager component.Manager
	router  router.Router
	manager component.Manager
}

func (p *manager) Init(opts ...Option) {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}

	if options.router != nil {
		p.router = options.router
	}

	if options.manager != nil {
		p.manager = options.manager
	}

	if p.router == nil {
		p.router = NewRoutes(options.logger)
	}
}

func (p *manager) CallComponent(ctx context.Context, msg message.Message) (interface{}, error) {

	cpt, err := p.manager.GetComponent(msg.Service())
	if err != nil {
		return nil, err
	} else if cpt == nil {
		return nil, fmt.Errorf("unknown component")
	}

	return cpt.Route(msg.Topic())(msg)
}

func (p *manager) CallServer(ctx context.Context, msg message.Message) (interface{}, error) {

	nodes, err := p.router.GetServiceNodes(router.ReadService(msg.Service()))
	if err != nil {
		return nil, err
	}

	nm, err := node.NewWithNodes(node.NodeTypeConsistent, msg.Service().FullRegistry(), nodes)
	if err != nil {
		return nil, err
	}

	var keys []string

	node, ok := nm.NodeFor(keys...)
	if !ok {
		return nil, fmt.Errorf("not found service node")
	}

	var rep interface{}
	switch node.Metadata["protocol"] {
	case service.Protocol_GRPC:
		fallthrough
	default:
		c := grpc.NewClient()

		// todo options
		req := c.NewRequest(msg.Service(), node.Value, msg.GetPayload().GetBody())
		ctx := context.Background()

		err := c.Call(ctx, req, rep)
		if err != nil {
			return nil, err
		}
	}

	return rep, nil
}

func (p *manager) Start() error {

	for _, cpt := range p.manager.ListComponents() {
		fmt.Println(cpt.Name, cpt.Component)
		if err := cpt.Component.Start(); err != nil {
			return err
		}
	}

	return p.router.Start()
}

func (p *manager) Stop() error {

	for _, cpt := range p.manager.ListComponents() {
		if err := cpt.Component.Stop(); err != nil {
			return err
		}
	}
	return p.router.Stop()
}

func (p *manager) Router() router.Router {
	return p.router
}

func (p *manager) CompManager() component.Manager {
	return p.manager
}
