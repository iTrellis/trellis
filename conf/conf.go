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

package conf

import (
	"github.com/go-trellis/config"
)

// Project project config info
type Project struct {
	Services map[string]Service `yaml:"services"`

	Registry []Registry `yam:"registries"`
}

// Service service info
type Service struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`

	Registry string `yaml:"registry"`

	Options config.Options `yaml:"options"`
}

// Registry run configure
type Registry struct {
	Name     string         `yaml:"name"`
	Type     string         `yaml:"type"`
	Options  config.Options `yaml:"options"`
	Protocal Protocal       `yaml:"protocal"`
}

// Protocal transfer's protocal
type Protocal struct {
	Type    string `yaml:"type"`
	Address string `yaml:"address"`
}
