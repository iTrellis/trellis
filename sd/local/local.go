package local

import (
	"github.com/iTrellis/trellis/service/component"
	"github.com/iTrellis/trellis/service/registry"
)

type localSD struct {
	components map[string]component.Component

	newComponentFuncs map[string]component.NewComponentFunc

	componentNames []string
}

// NewRegistry new default local route manager
func NewRegistry(opts ...registry.Option) registry.Registry {
	return nil
}

// func (*localSD) DeregisterRegistry(string) error {
// 	return nil
// }

// func (*localSD) DeregisterService(string, *registry.Service) error {
// 	return nil
// }

// // RegisterComponent regist component function
// func (p *defCompManager) RegisterComponentFunc(service *service.Service, fn component.NewComponentFunc) {

// 	if service == nil || len(service.Name) == 0 {
// 		panic("component name is empty")
// 	}

// 	if fn == nil {
// 		panic("component fn is nil")
// 	}

// 	_, exist := p.newComponentFuncs[service.FullPath()]

// 	if exist {
// 		panic(fmt.Sprintf("component already registered: %s", service.FullPath()))
// 	}

// 	p.newComponentFuncs[service.FullPath()] = fn
// 	p.componentNames = append(p.componentNames, service.FullPath())
// }

// // ListComponents get components
// func (p *defCompManager) ListComponents() []component.Describe {

// 	var desc []component.Describe

// 	for _, name := range p.componentNames {
// 		cpt := p.components[name]

// 		desc = append(desc, component.Describe{
// 			Name: name,
// 			// RegisterFunc: runtime.FuncForPC(reflect.ValueOf(cpt).Pointer()).Name(),
// 			RegisterFunc: reflect.ValueOf(cpt).String(),
// 			Component:    cpt,
// 		})
// 	}

// 	return desc
// }

// // NewComponent new component
// func (p *defCompManager) NewComponent(service *service.Service, alias string, opts ...component.Option) (
// 	component.Component, error) {
// 	fn, ok := p.newComponentFuncs[service.FullPath()]

// 	if !ok {
// 		return nil, fmt.Errorf("component driver '%s' not exist", service.FullPath())
// 	}

// 	cpt, err := fn(alias, opts...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	p.components[service.FullPath()] = cpt

// 	return cpt, nil
// }

// // GetComponent get component
// func (p *defCompManager) GetComponent(s *service.Service) (cpt component.Component, err error) {
// 	cpt, ok := p.components[s.FullPath()]
// 	if !ok {
// 		return nil, errors.New("component is not exists")
// 	}
// 	return cpt, nil
// }

// type compResp struct {
// 	r   interface{}
// 	err error
// }

// // GetComponent get component
// func (p *defCompManager) Call(msg message.Message, opts ...component.CallOption) (interface{}, error) {

// 	cpt, err := p.GetComponent(msg.Service())
// 	if err != nil {
// 		return nil, err
// 	}

// 	options := component.CallOptions{}
// 	for _, o := range opts {
// 		o(&options)
// 	}

// 	if options.Timeout == 0 {
// 		options.Timeout = 10 * time.Second
// 	}

// 	h := cpt.Route(msg.Topic())
// 	if h == nil {
// 		return nil, errors.New("not found handler")
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), options.Timeout)
// 	defer cancel()

// 	ch := make(chan compResp)

// 	go func() {
// 		respH, err := h(msg)
// 		ch <- compResp{
// 			r:   respH,
// 			err: err,
// 		}
// 	}()

// 	select {
// 	case res := <-ch:
// 		// return res
// 		return res.r, res.err
// 	case <-ctx.Done():
// 		return nil, errors.New("component exceed timeout")
// 	}
// }
