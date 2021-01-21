/*
Copyright © 2020 Henry Huang <hhh@rutcode.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"github.com/iTrellis/trellis/sd/memory"
	"github.com/iTrellis/trellis/service"
	"github.com/iTrellis/trellis/service/component"
	"github.com/iTrellis/trellis/service/registry"
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
