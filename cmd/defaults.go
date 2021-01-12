package cmd

import (
	"github.com/go-trellis/trellis/sd/memory"
	"github.com/go-trellis/trellis/service"
	"github.com/go-trellis/trellis/service/component"
	"github.com/go-trellis/trellis/service/registry"
)

var (
	DefaultNewRegistryFuncs = map[service.RegisterType]registry.NewRegistryFunc{
		// sd.RegistryETCD:
		// sd.RegistryMDNS:
		service.RegisterType_memory: memory.NewRegistry,
	}

	DefaultNewComponentFuncs = map[service.Service]component.NewComponentFunc{}

	DefaultHiddenVersions = []string{"0", "0.0", "0.0.0", "v0", "v0.0", "v0.0.0"}
)
