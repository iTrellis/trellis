package registry

// // ServiceOptionFunc 配置对象函数定义
// type ServiceOptionFunc func(*defService)

// // ServiceBaseService 基础服务信息
// func ServiceBaseService(base *proto.BaseService) ServiceOptionFunc {
// 	return func(s *defService) {
// 		s.BaseService = base
// 	}
// }

// // ServiceRegistryType 注册方式
// func ServiceRegistryType(typ proto.RegistryType) ServiceOptionFunc {
// 	return func(s *defService) {
// 		s.RegistryType = typ
// 	}
// }

// // ServiceName 服务名称
// func ServiceName(name string) ServiceOptionFunc {
// 	return func(s *defService) {
// 		s.Name = name
// 	}
// }

// // ServiceVersion 服务版本
// func ServiceVersion(version string) ServiceOptionFunc {
// 	return func(s *defService) {
// 		s.Version = version
// 	}
// }

// // ServiceNodes 服务节点信息
// func ServiceNodes(nodes node.Manager) ServiceOptionFunc {
// 	return func(s *defService) {
// 		s.NodeManager = nodes
// 	}
// }

// // ServiceMetadata 配置元数据
// func ServiceMetadata(data map[string]string) ServiceOptionFunc {
// 	return func(s *defService) {
// 		s.Metadata = data
// 	}
// }

// // Service 注册的Service对象
// type Service interface {
// 	GetRegistryType() proto.RegistryType
// 	GetName() string
// 	GetVersion() string
// 	GetMetadata() map[string]string
// 	GetNodes() node.Manager
// }

// type defService struct {
// 	*proto.BaseService
// 	proto.RegistryType
// 	Metadata    map[string]string
// 	NodeManager node.Manager
// }

// // NewService 生成新服务对象
// func NewService(opts ...ServiceOptionFunc) (Service, error) {
// 	s := &defService{}

// 	for _, o := range opts {
// 		o(s)
// 	}

// 	if s.Name == "" {
// 		return nil, errors.New("service's name is nil")
// 	}

// 	return s, nil
// }

// func (p *defService) GetRegistryType() proto.RegistryType {
// 	return p.RegistryType
// }

// func (p *defService) GetName() string {
// 	return p.Name
// }

// func (p *defService) GetVersion() string {
// 	return p.Version
// }

// func (p *defService) GetMetadata() map[string]string {
// 	return p.Metadata
// }

// func (p *defService) GetNodes() node.Manager {
// 	return p.NodeManager
// }
