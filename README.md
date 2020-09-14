# trellis 

A micro service framework, which can build some services (component) into one server (like lego) 

![pic](Trellis.jpg)

## Build Project

* You should use [tbuilder](https://github.com/go-trellis/tbuilder) for building your project
* projects' configure, see examples.

## Own Services

* Step 1: You can write your own services implement Service

```go
// Service 服务对象
type Service interface {
	LifeCycle
	Handlers
}

// Handlers 函数路由器
type Handlers interface {
	Route(topic string) HandlerFunc
}

// LifeCycle server的生命周期
type LifeCycle interface {
	Start() error
	Stop() error
}
```

* Step 2: regist your services into serviceFuncs

```go
type NewServiceFunc func(opts ...OptionFunc) (Service, error)

service.RegistNewServiceFunc(name, version, newFunc)
```

## More Examples

### Project's build configure

* [http post config with origin building](examples/cmd/build.yaml)
* [inner http api config](examples/cmd/build_remote.yaml)
* [inner grpc api config](examples/cmd/build_grpc.yaml)

### Project's run configure

* [http post config](examples/cmd/run.yaml)
* [inner http api config](examples/cmd/run_remote.yaml)
* [inner grpc api config](examples/cmd/run_grpc.yaml)

## SOMETHING TODO

* regitry: dns | consul ...
* calling chain service.
* config service: such as: build、 run configures.
* monitor service: monitor logs, services' status, etc.
